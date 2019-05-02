package com.tms.app

import org.scalatra._
import org.json4s.{DefaultFormats, Formats}
import org.scalatra.json._

class AccountsController extends ScalatraServlet with JacksonJsonSupport {

  protected implicit lazy val jsonFormats: Formats = DefaultFormats

  before() {
    contentType = formats("json")
  }

  get("/") {
    val account = Account("luispalominot@hotmail.com", "Luis Palomino", "luismfao2")
    account
  }

}

case class Account(email: String, name: String, password: String)