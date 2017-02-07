package com.reactivearchitecturecookbook.summary

import akka.actor.{Actor, OneForOneStrategy, Props, SupervisorStrategy}
import cakesolutions.kafka.akka.KafkaConsumerActor
import cakesolutions.kafka.akka.KafkaConsumerActor.Subscribe
import cakesolutions.kafka.{KafkaConsumer, KafkaDeserializer, KafkaProducer, KafkaSerializer}
import com.reactivearchitecturecookbook.Envelope
import com.redis.RedisClientPool
import com.typesafe.config.Config
import org.apache.kafka.common.serialization.{StringDeserializer, StringSerializer}

object SummaryRoot {
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

  import scala.concurrent.duration._
  override def supervisorStrategy: SupervisorStrategy = OneForOneStrategy(maxNrOfRetries = 3, withinTimeRange = 3.seconds) {
    case _ ⇒ SupervisorStrategy.Restart
  }


  @scala.throws[Exception](classOf[Exception])
  override def preStart(): Unit = kafkaConsumerActor ! Subscribe.AutoPartitionWithManualOffset()

  override def receive: Receive = {

  }

}
