package com.reactivearchitecturecookbook

import com.trueaccord.scalapb.json.JsonFormat

/**
  * Created by janmachacek on 16/01/2017.
  */
object Demo extends App {
  // This is the JSON produced by our C++ program
  val json = """{"correlationId":"d8437bf8-4f97-46b9-8b67-9d941c0cbe70","payload":{"@type":"type.googleapis.com/com.reactivearchitecturecookbook.faceextract.v1m0.ExtractFace","mimeType":"image/png","content":"PHBuZy1ieXRlcz4="}}"""
  // This is the ASCII representation of the bytes produced by our C++ program
  val bytes = "0A2433653730376463392D646466662D343261362D396566642D356165613736616562313636226D0A51747970652E676F6F676C65617069732E636F6D2F636F6D2E7265616374697665617263686974656374757265636F6F6B626F6F6B2E66616365657874726163742E76316D302E457874726163744661636512180A09696D6167652F706E67120B3C706E672D62797465733E"
    .grouped(2)
    .map(x ⇒ Integer.parseInt(x, 16).toByte)
    .toArray

  // we can parse from the bytes
  show("bytes", Envelope.parseFrom(bytes))
  // just as easily as from JSON
  show("json", JsonFormat.fromJsonString(json)(Envelope))

  // a convenience-only function that displays the payload
  private def show(source: String, envelope: Envelope): Unit = {
    envelope.payload
      .filter(_.is[faceextract.v1m0.ExtractFace])
      .map(_.unpack[faceextract.v1m0.ExtractFace])
      .foreach(x ⇒ println(s"$source:\n$x"))
  }

}
