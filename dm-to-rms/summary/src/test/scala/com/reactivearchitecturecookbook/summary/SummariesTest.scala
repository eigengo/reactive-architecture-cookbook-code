package com.reactivearchitecturecookbook.summary

import cakesolutions.kafka.akka.{ConsumerRecords, Offsets}
import com.reactivearchitecturecookbook.Envelope
import org.apache.kafka.common.TopicPartition
import org.scalatest.{FlatSpec, Matchers}

class SummariesTest extends FlatSpec with Matchers {

  case class ConsumerRecordsProducer(offsetsMap: Map[TopicPartition, Long]) {

    //noinspection VariablePatternShadow
    def appending(pairs: (String, Int, String, Envelope)*): (ConsumerRecordsProducer, ConsumerRecords[String, Envelope]) = {
      val crm = pairs.foldLeft(Map.empty[TopicPartition, Seq[(Option[String], Envelope)]]) {
        case (result, (topic, partition, key, payload)) ⇒
          val tp = new TopicPartition(topic, partition)

          result + ((tp, result.get(tp).map(records ⇒ records.+:(Some(key), payload)).getOrElse(List((Some(key), payload)))))
      }
      val om = pairs.foldLeft(offsetsMap) { case (offsetsMap, (topic, partition, _, _)) ⇒
        val tp = new TopicPartition(topic, partition)
        offsetsMap + ((tp, offsetsMap.get(tp).map(_ + 1).getOrElse(1L)))
      }
      (copy(offsetsMap = om), ConsumerRecords.fromMap(crm))
    }

  }

  it should "collect outcomes and offsets" in {
    val (crp, crs1) = ConsumerRecordsProducer(Map.empty).appending(
      ("a", 0, "a", Envelope()), ("a", 0, "b", Envelope()), ("a", 0, "a", Envelope()),
      ("a", 0, "b", Envelope()), ("a", 0, "a", Envelope())
    )
    val (summaries1, outcomes1, offsets1) = Summaries.empty.appending(crs1.recordsList)
    outcomes1 shouldBe Map(("a", Outcome(2)))
    offsets1 shouldBe Offsets(Map((new TopicPartition("a", 0), 1L)))

    val (_, crs2) = crp.appending(("a", 0, "b", Envelope()))
    val (_, outcomes2, offsets2) = summaries1.appending(crs2.recordsList)
    outcomes2 shouldBe Map(("b", Outcome(2)))
    offsets2 shouldBe Offsets(Map.empty)
  }

}
