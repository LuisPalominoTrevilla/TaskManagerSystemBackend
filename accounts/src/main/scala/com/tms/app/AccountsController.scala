package com.tms.app

import org.scalatra._
import org.json4s.{DefaultFormats, Formats}
import org.scalatra.json._
import org.scalatra.CorsSupport

import com.tms.db.Persistence
import com.tms.models._

class AccountsController extends ScalatraServlet with JacksonJsonSupport with CorsSupport {

  protected implicit lazy val jsonFormats: Formats = DefaultFormats

  before() {
    contentType = formats("json")
  }

  options("/*"){
    response.setHeader("Access-Control-Allow-Origin", request.getHeader("Origin"))
  }

  get("/") {
    "API is alive"
  }

  get("/healthCheck/?") {
    "OK"
  }

  get("/accounts/?") {
    try {
      Persistence.getAll
    } catch {
      case e: Throwable => {
          status = 500
          ResponseError(500, e.getMessage)
        }
    }
  }

  get("/accounts/:email/?") {
    try {
      val accounts: List[Account] = Persistence.getAll
      accounts find {
        case Account(email, _, _) => email == params("email")
      } match { case Some(account) => account; case None => throw new Error("Account was not found") }
    } catch {
      case e: Throwable => {
          status = 500
          ResponseError(500, e.getMessage)
        }
    }
  }

  post("/accounts/?") {
    try {
      val newAccount: Account = parsedBody.extract[Account]
      Persistence.insertOne(newAccount)
    } catch {
      case e: Throwable => {
        status = 400
        ResponseError(400, e.getMessage)
      }
    }
  }

  delete("/accounts/?") {
    contentType = "text/plain"
    try {
      Persistence.deleteAll
      "OK"
    } catch {
      case e: Throwable => {
        contentType = formats("json")
        status = 400
        ResponseError(400, e.getMessage)
      }
    }
  }

  delete("/accounts/:email/?") {
    contentType = "text/plain"
    try {
      Persistence.deleteOne(params("email"))
      "OK"
    } catch {
      case e: Throwable => {
        contentType = formats("json")
        status = 400
        ResponseError(400, e.getMessage)
      }
    }
  }

  post("/accounts/login/?") {
    try {
      val credentials: Credentials = parsedBody.extract[Credentials]
      Persistence.getAll find {
        case Account(email, _, password) => email == credentials.email && password == credentials.password
      } match { case Some(account) => account; case None => throw new Error("Credentials not valid") }
    } catch {
      case e: Throwable => {
        status = 401
        ResponseError(401, e.getMessage)
      }
    }
  }
}