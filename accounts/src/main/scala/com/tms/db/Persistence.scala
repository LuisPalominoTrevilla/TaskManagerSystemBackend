package com.tms.db

import com.tms.models.Account
import java.nio.file.Paths

object Persistence {

    val filePath = Paths.get(System.getProperty("user.dir"), "src", "main", "resources", "accounts.txt").toString

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
        AWSBucket.saveFile(this.filePath)
        val fileContent: List[String] = FileModifier.getFileContents(this.filePath)
        if (validateDuplicated(fileContent, account.email)) {
            val newLine: String = account match { case Account(email, name, password) => "%s,%s,%s".format(email, name, password)}
            val file = FileModifier.insertIntoFile(this.filePath, newLine :: fileContent)
            AWSBucket.uploadFile(file)
            FileModifier.deleteFile(file)
        } else {
            throw new Error("Account already exists")
        }
        account
    }

    def deleteAll(): Unit = {
        AWSBucket.saveFile(this.filePath)
        val fileContent: List[String] = FileModifier.getFileContents(this.filePath)
        val file = FileModifier.insertIntoFile(this.filePath, Nil)
        AWSBucket.uploadFile(file)
        FileModifier.deleteFile(file)
    }

    def deleteOne(deleteMail: String): Unit = {
        AWSBucket.saveFile(this.filePath)
        val fileContent: List[String] = FileModifier.getFileContents(this.filePath)
        val origLength = fileContent.length
        val newContent: List[String] = (for {
            line <- fileContent
            if (line.split(",") match { case Array(mail, _, _) => mail != deleteMail; case default => false})
        } yield line)
        if (newContent.length < origLength) {
            val file = FileModifier.insertIntoFile(this.filePath, newContent)
            AWSBucket.uploadFile(file)
            FileModifier.deleteFile(file)
        } else {
            throw new Error("Nothing to delete")
        }
    }
}