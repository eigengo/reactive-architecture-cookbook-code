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
      if (events.size > 1) Complete(Outcome(2))
      else copy(events = envelope :: events)
    }
  }

  case class Complete(outcome: Outcome) extends Summary {
    override val isComplete: Boolean = true
    override def next(envelope: Envelope): Summary = this
  }
}
