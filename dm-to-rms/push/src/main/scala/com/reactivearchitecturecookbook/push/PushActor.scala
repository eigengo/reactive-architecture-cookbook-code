package com.reactivearchitecturecookbook.push

import java.nio.file.{Files, Paths}
import java.security.KeyFactory
import java.security.spec.PKCS8EncodedKeySpec

import akka.actor.{Actor, OneForOneStrategy, Props, SupervisorStrategy}
import akka.http.scaladsl._
import akka.http.scaladsl.model._
import akka.stream.ActorMaterializer
import cakesolutions.kafka._
import cakesolutions.kafka.akka.KafkaConsumerActor.{Confirm, Subscribe}
import cakesolutions.kafka.akka.{ConsumerRecords, KafkaConsumerActor}
import com.google.protobuf.ByteString
import com.nimbusds.jose.JWEDecrypter
import com.nimbusds.jose.crypto.RSADecrypter
import com.nimbusds.jwt.EncryptedJWT
import com.reactivearchitecturecookbook.Envelope
import com.reactivearchitecturecookbook.push.v1m0.PushRequest
import com.redis.RedisClientPool
import com.typesafe.config.Config
import org.apache.kafka.common.TopicPartition
import org.apache.kafka.common.serialization.{StringDeserializer, StringSerializer}

import scala.concurrent.Future
import scala.util.Try

object PushActor {
  private val extractor = ConsumerRecords.extractor[String, Envelope]
  private case object Sent

  def props(config: Config): Props = {
    val privateKeySpec = new PKCS8EncodedKeySpec(Files.readAllBytes(Paths.get(config.getString("keyPath"), "jwt_rsa")))
    val privateKey = KeyFactory.getInstance("RSA").generatePrivate(privateKeySpec)
    val rsaDecrypter = new RSADecrypter(privateKey)

    val consumerConf = KafkaConsumer.Conf(
      config.getConfig("kafka.consumer-config"),
      keyDeserializer = new StringDeserializer,
      valueDeserializer = KafkaDeserializer(Envelope.parseFrom)
    )
    val consumerActorConf = KafkaConsumerActor.Conf()
    val producerConf = KafkaProducer.Conf(
      config.getConfig("kafka.producer-config"),
      new StringSerializer,
      KafkaSerializer[Envelope](_.toByteArray)
    )

    Props(classOf[PushOutputActor], consumerConf, consumerActorConf, producerConf, rsaDecrypter)
  }
}

class PushActor(consumerConf: KafkaConsumer.Conf[String, Envelope],
                consumerActorConf: KafkaConsumerActor.Conf,
                producerConf: KafkaProducer.Conf[String, Envelope],
                redisClientPool: RedisClientPool,
                jwtDecrypter: JWEDecrypter) extends Actor {

  private[this] val kafkaConsumerActor = context.actorOf(
    KafkaConsumerActor.props(consumerConf = consumerConf, actorConf = consumerActorConf, downstreamActor = self)
  )
  private[this] val kafkaProducer = KafkaProducer(producerConf)

  import scala.concurrent.duration._

  implicit val materializer = ActorMaterializer()
  private val pool = Http(context.system).superPool[(TopicPartition, Long)]()

  override def supervisorStrategy: SupervisorStrategy = OneForOneStrategy(maxNrOfRetries = 3, withinTimeRange = 3.seconds) {
    case _ ⇒ SupervisorStrategy.Restart
  }

  @scala.throws[Exception](classOf[Exception])
  override def preStart(): Unit = {
    kafkaConsumerActor ! Subscribe.AutoPartition(Seq("vision-1", "vision-internal-1"))
  }

  override def receive: Receive = {
    case PushActor.extractor(consumerRecords) ⇒
      val out = consumerRecords.recordsList.flatMap { record ⇒
        (for {
          jwt ← Try(EncryptedJWT.parse(record.value().token))
          _ ← Try(jwt.decrypt(jwtDecrypter))
          uri ← Try(Uri(jwt.getJWTClaimsSet.getStringClaim("push-*")))
          entity = HttpEntity(ContentTypes.`application/json`, "entity here") //JsonFormat.toJsonString(record.value().payload))

          payload = PushRequest(uri.toString(), ByteString.copyFrom(entity.data.asByteBuffer))
          newEnvelope = record.value().withPayload(com.google.protobuf.any.Any.pack(payload))
        } yield (record.key(), newEnvelope)).map { case (k, v) ⇒
          kafkaProducer.send(KafkaProducerRecord("push-internal-1", k, v))
        }.toOption
      }

      Future.sequence(out).foreach(_ ⇒ kafkaConsumerActor ! Confirm(consumerRecords.offsets))
  }

}
