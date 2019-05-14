const fs = require("fs")
const path = require("path")
const express = require("express")
const pathToSwaggerUi = require("swagger-ui-dist").absolutePath()

const indexContent = fs.readFileSync(`${pathToSwaggerUi}/index.html`)
  .toString()
  .replace("https://petstore.swagger.io/v2/swagger.json", "http://localhost:3000/static/taskManagerDoc.json")

const app = express()

app.get("/", (req, res) => res.send(indexContent))
app.get("/index.html", (req, res) => res.send(indexContent)) // you need to do this since the line below serves `index.html` at both routes
app.use(express.static(pathToSwaggerUi))
app.use('/static', express.static(path.join(__dirname, '/docfiles')))

app.listen(3000)