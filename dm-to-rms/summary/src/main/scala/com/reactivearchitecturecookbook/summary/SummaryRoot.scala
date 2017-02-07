package com.reactivearchitecturecookbook.summary

import akka.actor.{Actor, ActorRef, OneForOneStrategy, Props, SupervisorStrategy}
import cakesolutions.kafka.akka.KafkaConsumerActor.Subscribe
import cakesolutions.kafka.akka.{ConsumerRecords, KafkaConsumerActor, Offsets}
import cakesolutions.kafka.{KafkaConsumer, KafkaDeserializer, KafkaProducer, KafkaSerializer}
import com.reactivearchitecturecookbook.Envelope
import com.redis.RedisClientPool
import com.typesafe.config.Config
import org.apache.kafka.common.TopicPartition
import org.apache.kafka.common.serialization.{StringDeserializer, StringSerializer}

object SummaryRoot {
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
    val redisClientPool = new RedisClientPool(config.getString("redis.host"), config.getInt("redis.port"))

    Props(classOf[SummaryRoot], consumerConf, consumerActorConf, producerConf, redisClientPool)
  }
}

class SummaryRoot(consumerConf: KafkaConsumer.Conf[String, Envelope],
                  consumerActorConf: KafkaConsumerActor.Conf,
                  producerConf: KafkaProducer.Conf[String, Envelope],
                  redisClientPool: RedisClientPool) extends Actor {

  private[this] val kafkaConsumerActor = context.actorOf(
    KafkaConsumerActor.props(consumerConf = consumerConf, actorConf = consumerActorConf, downstreamActor = self)
  )
  private[this] var summaries: Map[String, (ActorRef, Offsets)] = Map.empty
  private[this] var offsets: Offsets = Offsets(Map.empty)

  import scala.concurrent.duration._
  override def supervisorStrategy: SupervisorStrategy = OneForOneStrategy(maxNrOfRetries = 3, withinTimeRange = 3.seconds) {
    case _ ⇒ SupervisorStrategy.Restart
  }

  def assignedListener(partitions: List[TopicPartition]): Offsets = {
    redisClientPool.withClient { client ⇒
      val offsets = partitions.map { p ⇒ (p, 0L) }
      Offsets(offsets.toMap)
    }
  }

  def revokedListener(partitions: List[TopicPartition]): Unit = {
    // noop
  }

  @scala.throws[Exception](classOf[Exception])
  override def preStart(): Unit = {
    kafkaConsumerActor ! Subscribe.AutoPartitionWithManualOffset(Seq("vision-1"), assignedListener, revokedListener)
  }

  private def summaryActorFor(transactionId: String, startingOffsets: Offsets): ActorRef = {
    if (summaries.contains(transactionId)) summaries(transactionId)._1
    else {
      val summary = context.actorOf(Summary.props, name = transactionId)
      summaries = summaries + ((transactionId, (summary, startingOffsets)))
      summary
    }
  }

  override def receive: Receive = {
    case SummaryRoot.extractor(consumerRecords) ⇒
      consumerRecords.pairs.foreach {
        case (Some(transactionId), envelope) ⇒ summaryActorFor(transactionId, offsets) ! envelope
        case _ ⇒ context.system.log.error("Received Envelope without a valid transactionId key.")
      }

      offsets = consumerRecords.offsets
      kafkaConsumerActor ! KafkaConsumerActor.Confirm(consumerRecords.offsets)
    case Summary.Completed(state) ⇒
      summaries = summaries - sender().path.name
      context.stop(sender())
      context.system.log.info(s"Completed with state $state.")
  }

}
