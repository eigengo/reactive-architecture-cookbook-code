package com.reactivearchitecturecookbook.push

import java.text.SimpleDateFormat
import java.util.{Date, UUID}

import akka.actor.ActorSystem
import com.nimbusds.jwt.JWTClaimsSet
import com.typesafe.config.{ConfigFactory, ConfigResolveOptions}

object Main extends App {
  System.setProperty("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
  System.setProperty("REDIS_HOST", "localhost")
  System.setProperty("REDIS_PORT", "6379")

  val config = ConfigFactory.load("push.conf").resolve(ConfigResolveOptions.defaults())
  val system = ActorSystem(name = "push-1_0_0", config = config)

  val claimsSet = new JWTClaimsSet.Builder()
      .issuer("https://id.reactivearchitecturecookbook.com")
      .subject("client")
      .audience("https://ips.reactivearchitecturecookbook.com")
      .notBeforeTime(new Date())
      .issueTime(new Date())
      .expirationTime(new SimpleDateFormat("yyyyMMdd").parse("20180101"))
      .jwtID(UUID.randomUUID().toString)

      .build()
  println(claimsSet.toJSONObject)

}
