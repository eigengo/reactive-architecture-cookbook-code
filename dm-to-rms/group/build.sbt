scalaVersion := "2.12.1"

libraryDependencies += "com.trueaccord.scalapb" %% "scalapb-json4s" % "0.1.6"

libraryDependencies += "com.trueaccord.scalapb" %% "scalapb-runtime" % com.trueaccord.scalapb.compiler.Version.scalapbVersion % "protobuf"

libraryDependencies += "com.nimbusds" % "nimbus-jose-jwt" % "4.34.1"

libraryDependencies += "org.zeromq" % "jeromq" % "0.3.5"

libraryDependencies += "org.scalatest" %% "scalatest" % "3.0.1" % "test"

PB.includePaths in Compile ++= Seq(file("../protocol"))

PB.protoSources in Compile := Seq(file("../protocol"))

PB.targets in Compile := Seq(
  scalapb.gen(flatPackage = true) -> (sourceManaged in Compile).value
)
