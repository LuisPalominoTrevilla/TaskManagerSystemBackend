# Accounts Microservice #

## Build & Run ##

```sh
$ sbt
> jetty:start
> browse
```

## Automatic Code Reloading ##

```sh
$ sbt
> ~;jetty:stop;jetty:start
> browse
```

If `browse` doesn't launch your browser, manually open [http://localhost:4000/](http://localhost:4000/) in your browser.
