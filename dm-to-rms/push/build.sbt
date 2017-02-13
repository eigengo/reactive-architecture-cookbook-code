scalaVersion := "2.12.1"

resolvers += Resolver.bintrayRepo("cakesolutions", "maven")

lazy val akkaVersion = "2.4.16"
lazy val akkaHttpVersion = "10.0.3"
lazy val scalaKafkaClientVersion = "0.10.1.2"

libraryDependencies ++= Seq(
  "com.typesafe.akka" %% "akka-actor" % akkaVersion,
  "com.typesafe.akka" %% "akka-testkit" % akkaVersion,
  "com.typesafe.akka" %% "akka-http" % akkaHttpVersion,
  "net.cakesolutions" %% "scala-kafka-client" % scalaKafkaClientVersion,
  "net.cakesolutions" %% "scala-kafka-client-akka" % scalaKafkaClientVersion,
  "com.trueaccord.scalapb" %% "scalapb-json4s" % "0.1.6",
  "com.trueaccord.scalapb" %% "scalapb-runtime" % com.trueaccord.scalapb.compiler.Version.scalapbVersion % "protobuf",
  "com.nimbusds" % "nimbus-jose-jwt" % "4.34.1",
  "org.scalatest" %% "scalatest" % "3.0.1" % "test",
  "net.debasishg" %% "redisclient" % "3.3"
)

PB.includePaths in Compile ++= Seq(file("../protocol"))

PB.protoSources in Compile := Seq(file("../protocol"))

PB.targets in Compile := Seq(
  scalapb.gen(flatPackage = true) -> (sourceManaged in Compile).value
)
