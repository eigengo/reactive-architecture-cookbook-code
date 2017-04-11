package com.reactivearchitecturecookbook.summary

import com.reactivearchitecturecookbook.Envelope

sealed trait Summary {
  def next(envelope: Envelope): Summary
  def isComplete: Boolean
}

object Summary {
  def apply(envelope: Envelope): Summary = {
    // here be some logic, but we say that one envelope never makes a complete summary
    Incomplete(List(envelope))
  }

  case class Incomplete private(events: List[Envelope]) extends Summary {
    override val isComplete: Boolean = false

    override def next(envelope: Envelope): Summary = {
      if (events.size < 2) copy(events = envelope :: events)
      else {
        val out = Envelope(events.head.correlationId, events.head.token, None)
        Complete(out)
      }
    }
  }

  case class Complete(outcome: Envelope) extends Summary {
    override val isComplete: Boolean = true
    override def next(envelope: Envelope): Summary = this
  }
}
