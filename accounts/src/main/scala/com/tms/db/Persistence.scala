package com.tms.db

import com.tms.models.Account

object Persistence {

    def getAll(): List[Account] = {
        def processAccounts(accounts: List[String], result: List[Account]): List[Account] = {
            accounts match {
                case account :: tail => {
                    val fields = account.split(",")
                    if (fields.length != 3) processAccounts(tail, result)
                    else {
                        val newAccount = Account(fields(0), fields(1), fields(2))
                        processAccounts(tail, newAccount :: result)
                    }
                }
                case Nil => result
            }
        }
        val rawAccounts: List[String] = AWSBucket.getFile()
        processAccounts(rawAccounts, Nil)
    }
}