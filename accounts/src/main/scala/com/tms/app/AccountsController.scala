package com.tms.app

import org.scalatra._

class AccountsController extends ScalatraServlet {

  get("/") {
    "Hola Mundo"
  }

}
