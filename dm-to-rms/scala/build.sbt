scalaVersion := "2.12.1"

libraryDependencies += "com.trueaccord.scalapb" %% "scalapb-json4s" % "0.1.6"

libraryDependencies += "com.trueaccord.scalapb" %% "scalapb-runtime" % com.trueaccord.scalapb.compiler.Version.scalapbVersion % "protobuf"

PB.includePaths in Compile ++= Seq(file("../protocol"))

PB.protoSources in Compile := Seq(file("../protocol"))

PB.targets in Compile := Seq(
  scalapb.gen(flatPackage = true) -> (sourceManaged in Compile).value
)
