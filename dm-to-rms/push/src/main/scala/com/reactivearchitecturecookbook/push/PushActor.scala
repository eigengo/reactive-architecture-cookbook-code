package com.reactivearchitecturecookbook.push

import akka.actor.{Actor, Kill, OneForOneStrategy, Props, SupervisorStrategy}
import cakesolutions.kafka._
import cakesolutions.kafka.akka.KafkaConsumerActor.Subscribe
import cakesolutions.kafka.akka.{ConsumerRecords, KafkaConsumerActor, Offsets}
import com.reactivearchitecturecookbook.Envelope
import com.redis.RedisClientPool
import com.typesafe.config.Config
import org.apache.kafka.common.TopicPartition
import org.apache.kafka.common.serialization.{StringDeserializer, StringSerializer}

import scala.concurrent.Future
import scala.util.{Failure, Success}

object PushActor {
  private val extractor = ConsumerRecords.extractor[String, Envelope]

  def props(config: Config): Props = {
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
    Props(classOf[PushActor], consumerConf, consumerActorConf, producerConf)
  }
}

class PushActor(consumerConf: KafkaConsumer.Conf[String, Envelope],
                consumerActorConf: KafkaConsumerActor.Conf,
                producerConf: KafkaProducer.Conf[String, Envelope],
                redisClientPool: RedisClientPool) extends Actor {

  private[this] val kafkaConsumerActor = context.actorOf(
    KafkaConsumerActor.props(consumerConf = consumerConf, actorConf = consumerActorConf, downstreamActor = self)
  )
  private[this] val kafkaProducer = KafkaProducer(producerConf)
  private[this] var lastFailedOffsetWrite: Long = 0

  import scala.concurrent.duration._

  private def revokedListener(partitions: List[TopicPartition]): Unit = {
    // noop
  }

  private def assignedListener(partitions: List[TopicPartition]): Offsets = {
    redisClientPool.withClient { client ⇒
      import com.redis.serialization.Parse.Implicits._
      val offsetsMap = partitions.map { tp ⇒ (tp, client.hget[Long](tp.topic(), tp.partition()).getOrElse(0L)) }.toMap
      Offsets(offsetsMap)
    }
  }

  override def supervisorStrategy: SupervisorStrategy = OneForOneStrategy(maxNrOfRetries = 3, withinTimeRange = 3.seconds) {
    case _ ⇒ SupervisorStrategy.Restart
  }


  @scala.throws[Exception](classOf[Exception])
  override def preStart(): Unit = {
    kafkaConsumerActor ! Subscribe.AutoPartitionWithManualOffset(
      Seq("vision-1", "vision-internal-1"),
      assignedListener,
      revokedListener
    )
  }

  override def receive: Receive = {
    case PushActor.extractor(consumerRecords) ⇒
      

  }

}
