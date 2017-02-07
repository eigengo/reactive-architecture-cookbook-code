package com.reactivearchitecturecookbook.summary

import akka.actor.ActorSystem
import com.typesafe.config.{ConfigFactory, ConfigResolveOptions}

object Main extends App {
  System.setProperty("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
  System.setProperty("REDIS_HOST", "localhost")
  System.setProperty("REDIS_PORT", "6379")

  val config = ConfigFactory.load("summary.conf").resolve(ConfigResolveOptions.defaults())
  val system = ActorSystem(name = "summary-1_0_0", config = config)
  system.actorOf(SummaryRoot.props(config.getConfig("app")))

//val jwe = new JWEObject(new JWEHeader(JWEAlgorithm.PBES2_HS512_A256KW, EncryptionMethod.A256CBC_HS512), new Payload("Hello, world"))
}
