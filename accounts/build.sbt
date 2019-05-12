val ScalatraVersion = "2.6.5"

organization := "com.tms"

name := "Accounts"

version := "0.0.1"

scalaVersion := "2.12.6"

resolvers += Classpaths.typesafeReleases

libraryDependencies ++= Seq(
  "org.scalatra" %% "scalatra" % ScalatraVersion,
  "org.scalatra" %% "scalatra-scalatest" % ScalatraVersion % "test",
  "ch.qos.logback" % "logback-classic" % "1.2.3" % "runtime",
  "org.eclipse.jetty" % "jetty-webapp" % "9.4.9.v20180320" % "container",
  "javax.servlet" % "javax.servlet-api" % "3.1.0" % "provided",
  "org.scalatra" %% "scalatra-json" % ScalatraVersion,
  "org.json4s"   %% "json4s-jackson" % "3.5.2",
  "com.amazonaws" % "aws-java-sdk" % "1.3.32",
  "commons-io" % "commons-io" % "2.6",
  "org.scalaj" %% "scalaj-http" % "2.4.1"
)

enablePlugins(ScalatraPlugin)

containerPort in Jetty := 4000