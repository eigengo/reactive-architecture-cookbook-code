package com.reactivearchitecturecookbook.push

import java.net.URL
import java.nio.file.{Files, Paths}
import java.security.KeyFactory
import java.security.spec.PKCS8EncodedKeySpec

import akka.actor.{Actor, OneForOneStrategy, Props, SupervisorStrategy}
import cakesolutions.kafka._
import cakesolutions.kafka.akka.KafkaConsumerActor.Subscribe
import cakesolutions.kafka.akka.{ConsumerRecords, KafkaConsumerActor}
import com.nimbusds.jose.JWEDecrypter
import com.nimbusds.jose.crypto.RSADecrypter
import com.nimbusds.jwt.EncryptedJWT
import com.reactivearchitecturecookbook.Envelope
import com.redis.RedisClientPool
import com.typesafe.config.Config
import org.apache.kafka.common.serialization.{StringDeserializer, StringSerializer}

import scala.util.Try

object PushActor {
  private val extractor = ConsumerRecords.extractor[String, Envelope]

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

    Props(classOf[PushActor], consumerConf, consumerActorConf, producerConf, rsaDecrypter)
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
  import scala.concurrent.duration._

  override def supervisorStrategy: SupervisorStrategy = OneForOneStrategy(maxNrOfRetries = 3, withinTimeRange = 3.seconds) {
    case _ ⇒ SupervisorStrategy.Restart
  }

  @scala.throws[Exception](classOf[Exception])
  override def preStart(): Unit = {
    kafkaConsumerActor ! Subscribe.AutoPartition(Seq("vision-1", "vision-internal-1"))
  }

  override def receive: Receive = {
    case PushActor.extractor(consumerRecords) ⇒
      val x = consumerRecords.recordsList.flatMap { record ⇒
        val encryptedJwt = EncryptedJWT.parse(record.value().token)
        for {
          _ ← Try(encryptedJwt.decrypt(jwtDecrypter)).toOption
          pushUrlString ← Option(encryptedJwt.getJWTClaimsSet.getStringClaim("push"))
          pushURL ← Try(new URL(pushUrlString)).toOption
        } yield (pushURL, record.key(), record.value())
      }.groupBy(_._1)
      x
  }

}
