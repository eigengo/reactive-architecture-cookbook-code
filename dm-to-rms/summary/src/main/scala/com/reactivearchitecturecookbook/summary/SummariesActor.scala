package com.reactivearchitecturecookbook.summary

import akka.actor.{Actor, OneForOneStrategy, Props, SupervisorStrategy}
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

object SummariesActor {
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

    Props(classOf[SummariesActor], consumerConf, consumerActorConf, producerConf, redisClientPool)
  }
}

class SummariesActor(consumerConf: KafkaConsumer.Conf[String, Envelope],
                     consumerActorConf: KafkaConsumerActor.Conf,
                     producerConf: KafkaProducer.Conf[String, Envelope],
                     redisClientPool: RedisClientPool) extends Actor {

  private[this] val kafkaConsumerActor = context.actorOf(
    KafkaConsumerActor.props(consumerConf = consumerConf, actorConf = consumerActorConf, downstreamActor = self)
  )
  private[this] val kafkaProducer = KafkaProducer(producerConf)
  private[this] var summaries: Summaries = Summaries.empty

  import scala.concurrent.duration._

  override def supervisorStrategy: SupervisorStrategy = OneForOneStrategy(maxNrOfRetries = 3, withinTimeRange = 3.seconds) {
    case _ ⇒ SupervisorStrategy.Restart
  }

  private def assignedListener(partitions: List[TopicPartition]): Offsets = {
    redisClientPool.withClient { client ⇒
      val offsets = partitions.map { p ⇒ (p, 0L) }
      Offsets(offsets.toMap)
    }
  }

  private def revokedListener(partitions: List[TopicPartition]): Unit = {
    // noop
  }

  private def persistOffsets(offsets: Offsets): Unit = {
    import context.dispatcher
    if (!offsets.isEmpty) Future {
      redisClientPool.withClient { client ⇒
        offsets.offsetsMap.foreach { case (tp, offset) ⇒ client.hset1(tp.topic(), tp.partition(), offset) }
        context.system.log.debug(s"Persisted latest offsets $offsets")
      }
    }
  }

  @scala.throws[Exception](classOf[Exception])
  override def preStart(): Unit = {
    kafkaConsumerActor ! Subscribe.AutoPartitionWithManualOffset(Seq("vision-1", "vision-internal-1"), assignedListener, revokedListener)
  }

  override def receive: Receive = {
    case SummariesActor.extractor(consumerRecords) ⇒
      val (newSummaries, outcomes, offsets) = summaries.withConsumerRecords(consumerRecords.recordsList)
      summaries = newSummaries
      kafkaConsumerActor ! KafkaConsumerActor.Confirm(consumerRecords.offsets)

      println(outcomes)
      println(offsets)
      if (outcomes.nonEmpty) {
        val sent = outcomes.map { case (transactionId, outcome) ⇒
          val out = Envelope()
          kafkaProducer.send(KafkaProducerRecord("summary-1", transactionId, out))
        }
        import context.dispatcher
        Future.sequence(sent).onComplete {
          case Success(recordMetadatas) ⇒
            context.system.log.info(s"Successfully sent $recordMetadatas for $offsets")
            persistOffsets(offsets)
          case Failure(ex) ⇒
            context.system.log.error(s"Did not send outcomes for $offsets because of $ex")
        }
      }

  }

}
