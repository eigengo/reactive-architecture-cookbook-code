package com.reactivearchitecturecookbook.push

import java.io.FileOutputStream
import java.nio.file.{Files, Paths}
import java.security.{KeyFactory, KeyPairGenerator}
import java.security.interfaces.RSAPublicKey
import java.security.spec.{PKCS8EncodedKeySpec, X509EncodedKeySpec}
import java.text.SimpleDateFormat
import java.util.{Date, UUID}

import akka.actor.ActorSystem
import com.nimbusds.jose.crypto.{RSADecrypter, RSAEncrypter}
import com.nimbusds.jose._
import com.nimbusds.jwt.{EncryptedJWT, JWTClaimsSet, SignedJWT}
import com.typesafe.config.{ConfigFactory, ConfigResolveOptions}

object Main extends App {
  System.setProperty("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
  System.setProperty("REDIS_HOST", "localhost")
  System.setProperty("REDIS_PORT", "6379")
  System.setProperty("KEY_PATH", "")

  val kp = KeyPairGenerator.getInstance("RSA").genKeyPair()
  def save(fn: String, contents: Array[Byte]): Unit = {
    val fos = new FileOutputStream(fn)
    fos.write(contents)
    fos.close()
  }
  save("jwt_rsa", kp.getPrivate.getEncoded)
  save("jwt_rsa.pub", kp.getPublic.getEncoded)

  val config = ConfigFactory.load("push.conf").resolve(ConfigResolveOptions.defaults())

  val privateKeySpec = new PKCS8EncodedKeySpec(Files.readAllBytes(Paths.get(config.getString("app.keyPath"), "jwt_rsa")))
  val publicKeySpec  = new X509EncodedKeySpec(Files.readAllBytes(Paths.get(config.getString("app.keyPath"), "jwt_rsa.pub")))
  val system = ActorSystem(name = "push-1_0_0", config = config)

  val jwtClaims = new JWTClaimsSet.Builder()
      .issuer("https://id.reactivearchitecturecookbook.com")
      .subject("client-account")
      .audience("https://ips.reactivearchitecturecookbook.com")
      .notBeforeTime(new Date())
      .issueTime(new Date())
      .expirationTime(new SimpleDateFormat("yyyyMMdd").parse("20180101"))
      .claim("push", "https://client.system.com")
      .jwtID(UUID.randomUUID().toString)
      .build()

  val header = new JWEHeader(JWEAlgorithm.RSA_OAEP, EncryptionMethod.A128GCM)
  val jwt = new EncryptedJWT(header, jwtClaims)
  val publicKey = KeyFactory.getInstance("RSA").generatePublic(publicKeySpec)
  val encrypter = new RSAEncrypter(publicKey.asInstanceOf[RSAPublicKey])
  jwt.encrypt(encrypter)
  val jwtString = jwt.serialize()
  println(jwtString)

  val jwt2 = EncryptedJWT.parse(jwtString)
  val privateKey = KeyFactory.getInstance("RSA").generatePrivate(privateKeySpec)
  val decrypter = new RSADecrypter(privateKey)
  jwt2.decrypt(decrypter)

  println(jwt2.getJWTClaimsSet.getStringClaim("push"))

}
