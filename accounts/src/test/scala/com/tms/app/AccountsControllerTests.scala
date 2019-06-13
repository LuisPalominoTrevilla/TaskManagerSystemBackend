package com.tms.app

import org.scalatra.test.scalatest._

class AccountsControllerTests extends ScalatraFunSuite {

  addServlet(classOf[AccountsController], "/*")

  test("GET / on AccountsController should return status 200") {
    get("/") {
      status should equal (200)
    }
  }

  test("GET /accounts/:fakeId shpuld return status 500") {
    get("/accounts/holazabmdiedl") {
      status should equal (500)
    }
  }

  test("POST /accounts should send an error 400 while attempting to create account") {
    post("/accounts") {
      status should equal (400)

    }
  }

  test("POST /accounts should create a new account") {
    post("/accounts/login") {
      status should equal (401)
    }
  }

  test("DELETE /accounts/fakeid should fail to delete a non-existent account") {
    delete("/accounts/34") {
      status should equal (400)

    }
  }

}
