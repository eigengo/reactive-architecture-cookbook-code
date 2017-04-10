package com.reactivearchitecturecookbook.push

import java.nio.ByteBuffer
import java.nio.file.{Files, Paths}
import java.security.KeyFactory
import java.security.spec.PKCS8EncodedKeySpec
import java.util.UUID

import akka.actor.ActorSystem
import akka.http.scaladsl.Http
import akka.http.scaladsl.model._
import akka.http.scaladsl.server.{RejectionError, UnsupportedRequestContentTypeRejection}
import akka.http.scaladsl.util.FastFuture
import akka.stream.ActorMaterializer
import cakesolutions.kafka.{KafkaProducer, KafkaProducerRecord, KafkaSerializer}
import com.google.protobuf.ByteString
import com.nimbusds.jose.crypto.RSADecrypter
import com.nimbusds.jwt.{EncryptedJWT, JWTClaimsSet}
import com.reactivearchitecturecookbook.Envelope
import com.reactivearchitecturecookbook.ingest.v1m0.IngestedImage
import com.typesafe.config.{ConfigFactory, ConfigResolveOptions}
import org.apache.kafka.common.serialization.StringSerializer

import scala.concurrent.{ExecutionContext, Future}
import scala.util.Try

object Main extends App with IngestRoute {
  System.setProperty("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
  System.setProperty("KEY_PATH", "")

  val config = ConfigFactory.load("ingest.conf").resolve(ConfigResolveOptions.defaults())
  implicit val system: ActorSystem = ActorSystem("ingest_1_0_0", config)
  implicit val materializer: ActorMaterializer = ActorMaterializer()

  val decrypter = {
    val privateKeySpec = new PKCS8EncodedKeySpec(Files.readAllBytes(Paths.get(config.getString("app.keyPath"), "jwt_rsa")))
    val privateKey = KeyFactory.getInstance("RSA").generatePrivate(privateKeySpec)
    new RSADecrypter(privateKey)
  }

  lazy val kafkaProducer = {
    val conf = KafkaProducer.Conf(
      config.getConfig("app.kafka.producer-config"),
      new StringSerializer,
      KafkaSerializer[Envelope](_.toByteArray)
    )
    KafkaProducer(conf)
  }

  override def authorizeAndExtract(token: String): Try[JWTClaimsSet] = {
    for {
      jwt ← Try(EncryptedJWT.parse(token))
      _ ← Try(jwt.decrypt(decrypter))
      cs = jwt.getJWTClaimsSet
      if cs.getStringClaim("cid") != null
    } yield cs
  }

  override def extractIngestImage(request: HttpRequest)(implicit ec: ExecutionContext): Future[IngestedImage] = {
    def validate(contentType: ContentType)(entity: HttpEntity.Strict): Future[IngestedImage] = {
      val image = IngestedImage(contentType.mediaType.value, ByteString.copyFrom(entity.data.asByteBuffer))
      // TODO: validate
      FastFuture.successful(image)
    }

    import scala.concurrent.duration._
    val supported = Set(ContentTypeRange(mediaRange = MediaRanges.`image/*`))

    request.header[akka.http.scaladsl.model.headers.`Content-Type`]
      .filter(ct ⇒ supported.exists(_.matches(ct.contentType)))
      .map { ct ⇒
        request.entity
          .withSizeLimit(5L * 1024L * 1024L)
          .toStrict(300.millis)
          .flatMap(validate(ct.contentType))
      }.getOrElse(FastFuture.failed(RejectionError(UnsupportedRequestContentTypeRejection(supported))))
  }

  override def publishIngestedImage(claimsSet: JWTClaimsSet, token: String, transactionId: String)(ingestedImage: IngestedImage)(implicit ec: ExecutionContext): Future[Unit] = {
    val clientId = claimsSet.getStringClaim("cid")

    val payload = com.google.protobuf.any.Any.pack(ingestedImage)
    val envelope = Envelope(correlationIds = Seq(UUID.randomUUID().toString), token, Some(payload))

    kafkaProducer.send(KafkaProducerRecord(clientId, envelope)).map(println)
  }

  import system.dispatcher

  Http().bindAndHandle(ingestRoute, "localhost", 9000)
}
