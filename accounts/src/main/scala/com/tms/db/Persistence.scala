package com.tms.db

import scala.io.Source
import com.amazonaws.auth.BasicAWSCredentials
import com.amazonaws.services.s3.AmazonS3Client
import com.amazonaws.AmazonClientException
import com.amazonaws.AmazonServiceException
import java.io.BufferedReader
import java.io.InputStreamReader

object Persistence {

    val BUCKET_NAME = "tms-itesm"
    val FILE_PATH = "/"
    val DATABASE_FILE = "accounts.txt"
    val AWS_ACCESS_KEY = "AKIAIFAZARV4736VE7WQ"
    val AWS_SECRET_KEY = "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"

    def readFile(reader: BufferedReader): List[String] = {
        val line = reader.readLine
        if (line != null) line :: readFile(reader)
        else Nil
    }
    
    def getAll(): List[String] = {
        try {
            val awsCredentials = new BasicAWSCredentials(AWS_ACCESS_KEY, AWS_SECRET_KEY)
            val amazonS3Client = new AmazonS3Client(awsCredentials)

            val obj = amazonS3Client.getObject(BUCKET_NAME, DATABASE_FILE)
            val reader = new BufferedReader(new InputStreamReader(obj.getObjectContent()))
            val lines: List[String] = readFile(reader)
            lines
        } catch {
            case ase: AmazonServiceException => throw new Exception("An error occurred")
            case ace: AmazonClientException => throw new Exception("An error occurred")
        }
        /* val bufferedSource = Source.fromFile("https://s3.us-east-2.amazonaws.com/task-management-system-itesm/accounts.txt")
        val documents: List[String] = {
            for {
                line <- bufferedSource.getLines.toList
            } yield line
        }
        bufferedSource.close
        documents */
    }
}