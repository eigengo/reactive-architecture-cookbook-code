package com.reactivearchitecturecookbook.summary

import akka.actor.{Actor, Props}
import com.reactivearchitecturecookbook.Envelope
import com.reactivearchitecturecookbook.faceextract.v1m0.ExtractedFace

object Summary {
  type State = Int

  case class Completed(state: State)

  val props: Props = Props[Summary]

}

class Summary extends Actor {
  import Summary._
  private var counter: Int = 0

  override def receive: Receive = {
    case Envelope(correlationId, Some(payload)) ⇒
      if (payload.is[ExtractedFace]) counter += 1
      if (counter == 2) context.parent ! Completed(counter)
  }
}