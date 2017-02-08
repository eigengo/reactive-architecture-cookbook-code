package com.reactivearchitecturecookbook.summary

import cakesolutions.kafka.akka.Offsets
import com.reactivearchitecturecookbook.Envelope
import org.apache.kafka.clients.consumer.ConsumerRecord
import org.apache.kafka.common.TopicPartition

object Summaries {
  val empty = Summaries(Map.empty)

  object SummaryWithOffsets {
    def apply(summary: Summary, topic: String, partition: Int, offset: Long): SummaryWithOffsets = {
      SummaryWithOffsets(summary, Map(new TopicPartition(topic, partition) → offset))
    }
  }

  case class SummaryWithOffsets private(summary: Summary, topicPartitionOffsets: Map[TopicPartition, Long]) {
    def next(envelope: Envelope, topic: String, partition: Int, offset: Long): SummaryWithOffsets = {
      val tp = new TopicPartition(topic, partition)
      copy(summary = summary.next(envelope), topicPartitionOffsets = topicPartitionOffsets + (tp → offset))
    }
  }

}

case class Summaries private(summaries: Map[String, Summaries.SummaryWithOffsets]) {

  import Summaries._

  def withConsumerRecords(consumerRecords: List[ConsumerRecord[String, Envelope]]): (Summaries, Map[String, Outcome], Offsets) = {
    //noinspection VariablePatternShadow
    val (added, updated) = consumerRecords.foldLeft((Map.empty[String, SummaryWithOffsets], Map.empty[String, SummaryWithOffsets])) {
      case ((added, updated), cr) ⇒
        val transactionId = cr.key()
        val envelope = cr.value()
        val topic = cr.topic()
        val partition = cr.partition()
        val offset = cr.offset()
        summaries.get(transactionId) match {
          case Some(existing) ⇒ (added, updated + ((transactionId, existing.next(envelope, topic, partition, offset))))
          case None ⇒ (added + ((transactionId, SummaryWithOffsets(Summary(envelope), topic, partition, offset))), updated)
        }
    }

    val completed = updated.filter { case (_, SummaryWithOffsets(s, _)) ⇒ s.isComplete }
    val newSummaries = summaries ++ added ++ updated -- completed.keys

    val offsets = {
      val nsv = newSummaries.values.withFilter(!_.summary.isComplete)
      val topicPartitions = nsv.flatMap(_.topicPartitionOffsets.keys).toSet
      val offsetsMap = topicPartitions.foldLeft(Map.empty[TopicPartition, Long]) { case (r, tp) ⇒
        val offsets = nsv.flatMap(_.topicPartitionOffsets.get(tp))
        if (offsets.isEmpty) r
        else {
          val latestOffset = r.get(tp).map(offset ⇒ math.max(offset, offsets.min)).getOrElse(offsets.min)
          r + ((tp, latestOffset))
        }
      }
      Offsets(offsetsMap)
    }

    val outcomes = completed.map { case (k, SummaryWithOffsets(Summary.Complete(outcome), _)) ⇒ (k, outcome) }
    (copy(summaries = newSummaries), outcomes, offsets)
  }

}