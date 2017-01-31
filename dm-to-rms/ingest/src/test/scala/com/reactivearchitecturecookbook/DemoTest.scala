package com.reactivearchitecturecookbook

import org.scalatest.{FlatSpec, Matchers}

class DemoTest extends FlatSpec with Matchers {

  it should "Translate" in {
    new Demo().foo(2) shouldBe "2"

    //new Demo().foo(4) shouldBe "fff"
  }

}
