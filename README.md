# Book Library

A simple RestFUL API for a book Library.

## Packages Used
- [rs/zerolog](https://github.com/rs/zerolog)
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [cznic/sqlite](https://gitlab.com/cznic/sqlite)
- [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
- [swaggo/swag](https://github.com/swaggo/swag)


## Usage
To run:
```shell
go run main.go
```
To run on a different port:
```shell
go run main.go --addr :<PORT>
```

## Logging
The application using zerolog module to support log levels. Default log level is set to error.
To run with debug logging level:
```shell
go run main.go --debug
```
The api requests are log to api.log and stdout. The custom format log is:
```shell
[Mon, 13 Feb 2023 14:01:25 MST] - ::1 "GET /v1/swagger/index.html HTTP/1.1 200 4.8455ms "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36" "
```

## API Documentation
Documentation generate with [swaggo/swag](https://github.com/swaggo/swag). 

Run application and browse to http://localhost:8080/v1/swagger/index.html

## Not Implemented
- Authentication
- TLS