package com.tms.db

import com.tms.models.Account
import java.nio.file.Paths

object Persistence {

    private def validateDuplicated(records: List[String], newRecord: String): Boolean = {
        records match {
            case record :: xs => if (record.split(",")(0) == newRecord) false else validateDuplicated(xs, newRecord)
            case Nil => true
        }
    }

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

    def insertOne(account: Account): Account = {
        val root = System.getProperty("user.dir")
        val filePath = Paths.get(root, "src", "main", "resources", "accounts.txt").toString
        AWSBucket.saveFile(filePath)
        val fileContent: List[String] = FileModifier.getFileContents(filePath)
        if (validateDuplicated(fileContent, account.email)) {
            val newLine: String = account match { case Account(email, name, password) => "%s,%s,%s".format(email, name, password)}
            val file = FileModifier.insertIntoFile(filePath, newLine :: fileContent)
            AWSBucket.uploadFile(file)
            FileModifier.deleteFile(file)
        } else {
            throw new Error("Account already exists")
        }
        account
    }

    def deleteAll(): Unit = {
        val root = System.getProperty("user.dir")
        val filePath = Paths.get(root, "src", "main", "resources", "accounts.txt").toString
        AWSBucket.saveFile(filePath)
        val fileContent: List[String] = FileModifier.getFileContents(filePath)
        val file = FileModifier.insertIntoFile(filePath, Nil)
        AWSBucket.uploadFile(file)
        FileModifier.deleteFile(file)
    }
}