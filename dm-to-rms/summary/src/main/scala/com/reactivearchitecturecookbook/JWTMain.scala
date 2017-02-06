package com.reactivearchitecturecookbook

import com.nimbusds.jose._

object JWTMain extends App {
  val jwe = new JWEObject(new JWEHeader(JWEAlgorithm.PBES2_HS512_A256KW, EncryptionMethod.A256CBC_HS512), new Payload("Hello, world"))
  
}
