# Accounts Microservice #

Accounts microservice that belongs to the task manager system

## Requirements

- [Java version 8](https://www.java.com/en/download/mac_download.jsp) is required. Run `javac -version` and `java -version`
- You will need [sbt](https://www.scala-sbt.org/). if you are using MacOS, you can get it with `brew install sbt`

## Build & Run

```sh
$ sbt
> jetty:start
> browse
```

## Automatic Code Reloading

```sh
$ sbt
> ~;jetty:stop;jetty:start
> browse
```

If `browse` doesn't launch your browser, manually open [http://localhost:4000/](http://localhost:4000/) in your browser.
