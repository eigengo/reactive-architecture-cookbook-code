package com.reactivearchitecturecookbook.summary

import com.reactivearchitecturecookbook.Envelope
import org.scalatest.{FlatSpec, Matchers}

class SummaryTest extends FlatSpec with Matchers {
  it should "work" in {
    val i1@Summary.Incomplete(_) = Summary(Envelope())
    val i2@Summary.Incomplete(_) = i1.next(Envelope())
    val Summary.Complete(envelope) = i2.next(Envelope())
  }
}
