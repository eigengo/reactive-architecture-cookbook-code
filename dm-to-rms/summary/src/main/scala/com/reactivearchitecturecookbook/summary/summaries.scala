package com.reactivearchitecturecookbook.summary

import cakesolutions.kafka.akka.Offsets
import com.reactivearchitecturecookbook.Envelope
import org.apache.kafka.clients.consumer.ConsumerRecord
import org.apache.kafka.common.TopicPartition

object Summaries {
  val empty = Summaries(Map.empty)

  object InternalSummary {
    def apply(summary: Summary, topic: String, partition: Int, offset: Long): InternalSummary = {
      InternalSummary(summary, Map(new TopicPartition(topic, partition) → (offset, offset)))
    }
  }

  case class InternalSummary private(summary: Summary, topicParittionOffsets: Map[TopicPartition, (Long, Long)]) {
    def next(envelope: Envelope, topic: String, partition: Int, offset: Long): InternalSummary = {
      val tp = new TopicPartition(topic, partition)
      val (first, _) = topicParittionOffsets(tp)
      copy(summary = summary.next(envelope), topicParittionOffsets = topicParittionOffsets + (tp → (first, offset)))
    }
  }

}

case class Summaries private(summaries: Map[String, Summaries.InternalSummary]) {

  import Summaries._

  private def completeOffsets(summaries: Map[String, Summaries.InternalSummary]): Offsets = {
    def computeOffsets(summaries: Iterable[InternalSummary]): Offsets = {
      val topicPartitions = summaries.flatMap(_.topicParittionOffsets.keys)
      val topicPartitionsOffsetRanges = topicPartitions.map { topicPartition ⇒
        (topicPartition, summaries.flatMap { summary ⇒ summary.topicParittionOffsets.get(topicPartition).map { case (start, end) ⇒ (summary.summary.isComplete, start, end) } })
      }

      val x= topicPartitionsOffsetRanges.flatMap { case (tp, offsetRanges) ⇒
        val sortedOffsetRanges = offsetRanges.toList.sortWith { case ((_, s1, _), (_, s2, _)) ⇒ s1 < s2 }
        //noinspection VariablePatternShadow
        val (minIncomplete, maxComplete) = sortedOffsetRanges.foldLeft((Long.MaxValue, 0L)) { case (r@(minIncomplete, maxComplete), (c, s, e)) ⇒
            if (c && maxComplete == 0) (minIncomplete, s)
            else if (!c && minIncomplete == Long.MaxValue) (e, maxComplete)
            else r
        }
        if (maxComplete <= minIncomplete) Some(tp, maxComplete) else None
      }
      Offsets(x.toMap)
    }

    computeOffsets(summaries.values)
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
