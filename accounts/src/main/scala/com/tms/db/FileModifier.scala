package com.tms.db

import java.io.File
import java.io.PrintWriter
import scala.io.Source

object FileModifier {

    def getFileContents(path: String): List[String] = {
        Source.fromFile(path).getLines.toList
    }

    def insertIntoFile(destination: String, lines: List[String]): File = {
        val file = new File(destination)
        val writer = new PrintWriter(file)
        val len = lines.length
        lines.zipWithIndex.foreach {
            case (line, idx) => {
                if (idx < len - 1) writer.println(line)
                else writer.print(line)
            }
        }
        writer.close
        file
    }

    def deleteFile(fileToDelete: File): Unit = fileToDelete.delete
}