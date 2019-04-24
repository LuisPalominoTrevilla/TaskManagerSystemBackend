package com.tms.app

import org.scalatra.test.scalatest._

class AccountsControllerTests extends ScalatraFunSuite {

  addServlet(classOf[AccountsController], "/*")

  test("GET / on AccountsController should return status 200") {
    get("/") {
      status should equal (200)
    }
  }

}
