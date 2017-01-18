package com.reactivearchitecturecookbook

import java.util.UUID

import com.google.protobuf.ByteString
import org.zeromq.{ZContext, ZMQ}

import scala.util.Try

object Main extends App {
  val context = new ZContext(2)
  val out = context.createSocket(ZMQ.PUSH)
  val in = context.createSocket(ZMQ.SUB)
  out.setSendTimeOut(1000)
  in.setReceiveTimeOut(1000)

  out.bind("tcp://*:5555")
  in.connect("tcp://localhost:5556")
  in.subscribe(Array.empty)

  while (true) {
    Try {
      val extractFace = faceextract.v1m0.ExtractFace("image/png", ByteString.EMPTY)
      val payload = com.google.protobuf.any.Any.pack(extractFace)
      val x = out.send(Envelope(UUID.randomUUID().toString, Some(payload)).toByteArray)
      println(s"Sent $x")
      Envelope.parseFrom(in.recv()).toString
    }
      .recover { case _ ⇒ "Failure" }
      .foreach { println }

  }
}
