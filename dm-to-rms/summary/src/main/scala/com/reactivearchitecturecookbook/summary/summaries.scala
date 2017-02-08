package com.reactivearchitecturecookbook.summary

import cakesolutions.kafka.akka.Offsets
import com.reactivearchitecturecookbook.Envelope
import org.apache.kafka.clients.consumer.ConsumerRecord
import org.apache.kafka.common.TopicPartition

object Summaries {
  val empty = Summaries(Map.empty)

  case class TopicParittionOffset(topic: String, partition: Int, offset: Long)

  object InternalSummary {
    def apply(summary: Summary, topic: String, partition: Int, offset: Long): InternalSummary = {
      InternalSummary(summary, List(TopicParittionOffset(topic, partition, offset)))
    }
  }

  case class InternalSummary private(summary: Summary, topicParittionOffsets: List[TopicParittionOffset]) {
    def next(envelope: Envelope, topic: String, partition: Int, offset: Long): InternalSummary = {
      val tpes = if (!topicParittionOffsets.exists(tpe ⇒ tpe.partition == partition && tpe.topic == topic)) {
        TopicParittionOffset(topic, partition, offset) :: topicParittionOffsets
      } else topicParittionOffsets
      copy(summary = summary.next(envelope), topicParittionOffsets = topicParittionOffsets)
    }
  }

}

case class Summaries private(summaries: Map[String, Summaries.InternalSummary]) {

  import Summaries._

  private def completeOffsets(summaries: Map[String, Summaries.InternalSummary]): Offsets = {
    val (cs, ic) = summaries.values.partition(_.summary.isComplete)
    val incompleteTpos = ic.flatMap(_.topicParittionOffsets)

    val complete = cs
        .flatMap(_.topicParittionOffsets)
        .groupBy(tpo ⇒ new TopicPartition(tpo.topic, tpo.partition))
        .mapValues(_.map(_.offset).min)
        .filterNot { case (tp, offset) ⇒ incompleteTpos.exists(x ⇒ x.topic == tp.topic() && x.partition == tp.partition() && x.offset > offset) }
//
//    if (complete.get(new TopicPartition("vision-1", 0)).contains(3503)) {
//      println("sfsfsdfs")
//    }

    Offsets(complete)
  }

  private def completeSummaries(summaries: Map[String, Summaries.InternalSummary]): Map[String, Outcome] = {
    summaries.flatMap {
      case (k, InternalSummary(CompleteSummary(outcome), _)) ⇒ Some(k, outcome)
      case _ ⇒ None
    }
  }

  def withConsumerRecords(consumerRecords: List[ConsumerRecord[String, Envelope]]): (Summaries, Map[String, Outcome], Offsets) = {
    val newSummaries = consumerRecords.foldLeft(summaries) { case (result, cr) ⇒
      val transactionId = cr.key()
      val envelope = cr.value()
      val topic = cr.topic()
      val partition = cr.partition()
      val offset = cr.offset()
      val internalSummary = result
        .get(transactionId)
        .map(_.next(envelope, topic, partition, offset))
        .getOrElse(InternalSummary(IncompleteSummary(envelope), topic, partition, offset))
      result + (transactionId → internalSummary)
    }
    val outcomes = completeSummaries(newSummaries)
    val offsets = completeOffsets(newSummaries)
    (copy(summaries = newSummaries -- outcomes.keys), outcomes, offsets)
  }

}
