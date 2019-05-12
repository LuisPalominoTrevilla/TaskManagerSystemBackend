## Habits Microservice #

Habits microservice that belongs to the task manager system

## Requirements

- [go 1.11](https://golang.org/)
- [mongodb](https://www.mongodb.com/)

## Run

- Use the following command

    `$ go run src/main.go`

- Browse `localhost:4001`

## Notes

- In case you get the following error:
    `cannot find module providing package golang.org/x/text/unicode/norm`
    You can get rid of it by simply running the following command:
    `go clean -modcache`
