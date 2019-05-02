package com.tms.app

import org.scalatra._
import org.json4s.{DefaultFormats, Formats}
import org.scalatra.json._
import com.tms.db.Persistence

class AccountsController extends ScalatraServlet with JacksonJsonSupport {

  protected implicit lazy val jsonFormats: Formats = DefaultFormats

  before() {
    contentType = formats("json")
  }

  get("/") {
    val account = Account("luispalominot@hotmail.com", "Luis Palomino", "luismfao2")
    account
  }

  get("/accounts") {
    Persistence.getAll()
  }

}

case class Account(email: String, name: String, password: String)