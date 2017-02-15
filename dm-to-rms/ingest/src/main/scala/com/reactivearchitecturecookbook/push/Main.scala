package com.reactivearchitecturecookbook.push

import java.util.UUID

import akka.actor.ActorSystem
import akka.http.scaladsl.Http
import akka.http.scaladsl.model.HttpRequest
import akka.stream.ActorMaterializer
import cakesolutions.kafka.{KafkaProducer, KafkaProducerRecord, KafkaSerializer}
import com.reactivearchitecturecookbook.Envelope
import com.reactivearchitecturecookbook.ingest.v1m0.IngestedImage
import com.typesafe.config.{ConfigFactory, ConfigResolveOptions}
import org.apache.kafka.common.serialization.StringSerializer

import scala.concurrent.{ExecutionContext, Future}

object Main extends App with IngestRoute {
  System.setProperty("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
  System.setProperty("KEY_PATH", "")

  val config = ConfigFactory.load("ingest.conf").resolve(ConfigResolveOptions.defaults())
  implicit val system: ActorSystem = ActorSystem("ingest_1_0_0", config)
  implicit val materializer: ActorMaterializer = ActorMaterializer()

  lazy val kafkaProducer = {
    val conf = KafkaProducer.Conf(
      config.getConfig("app.kafka.producer-config"),
      new StringSerializer,
      KafkaSerializer[Envelope](_.toByteArray)
    )
    KafkaProducer(conf)
  }

  override def authorizeJwt(token: String): Boolean = true

  override def extractIngestImage(request: HttpRequest)(implicit ec: ExecutionContext): Future[IngestedImage] = {
    Future(IngestedImage())
  }

  override def publishIngestedImage(token: String, transactionId: String)(ingestedImage: IngestedImage)(implicit ec: ExecutionContext): Future[Unit] = {
    val payload = com.google.protobuf.any.Any.pack(ingestedImage)
    val envelope = Envelope(correlationIds = Seq(UUID.randomUUID().toString), token, Some(payload))

    kafkaProducer.send(KafkaProducerRecord(transactionId, envelope)).map(println)
  }

  import system.dispatcher
  Http().bindAndHandle(ingestRoute, "localhost", 9000)
}
