package com.tms.app

import org.scalatra._
import org.json4s.{DefaultFormats, Formats}
import org.scalatra.json._

import com.tms.db.Persistence
import com.tms.models.{Account, Error}

class AccountsController extends ScalatraServlet with JacksonJsonSupport {

  protected implicit lazy val jsonFormats: Formats = DefaultFormats

  before() {
    contentType = formats("json")
  }

  get("/") {
    "API is alive"
  }

  get("/accounts") {
    try {
      Persistence.getAll
    } catch {
      case e: Throwable => {
          status = 500
          Error(500, e.getMessage)
        }
    }
  }

  post("/accounts") {
    try {
      val newAccount: Account = parsedBody.extract[Account]
      Persistence.insertOne(newAccount)
    } catch {
      case e: Throwable => {
        status = 400
        Error(400, e.getMessage)
      }
    }
  }
}