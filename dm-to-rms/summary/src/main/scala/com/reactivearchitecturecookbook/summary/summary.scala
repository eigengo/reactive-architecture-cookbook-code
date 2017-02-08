package com.reactivearchitecturecookbook.summary

import com.reactivearchitecturecookbook.Envelope

sealed trait Summary {
  def next(envelope: Envelope): Summary
  def isComplete: Boolean
}

object IncompleteSummary {
  def apply(envelope: Envelope): IncompleteSummary = IncompleteSummary(List(envelope))
}

case class IncompleteSummary private(events: List[Envelope]) extends Summary {
  override val isComplete: Boolean = false

  override def next(envelope: Envelope): Summary = {
    if (events.size > 2) CompleteSummary(Outcome(2))
    else copy(events = envelope :: events)
  }
}

case class Outcome(value: Int)

case class CompleteSummary(outcome: Outcome) extends Summary {
  override val isComplete: Boolean = true
  override def next(envelope: Envelope): Summary = this
}
