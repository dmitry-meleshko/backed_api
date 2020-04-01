# Backend - Patients

## Objective

Please demonstrate technical proficiency by creating a backend server application along with documentation and tests.

## Requirements

The backend stack must use technologies we use in-house:

* Golang
* PostgreSQL
* docker-compose

Implement the following REST APIs:

| Method | URL                             | Description                       |
|--------|---------------------------------|-----------------------------------|
| GET    | /api/v1/patients                | Get all patients                  |
| GET    | /api/v1/patients/:id            | Get one patient                   |
| POST   | /api/v1/patients                | Add one patient                   |

* The APIs must be auth-protected.
* The request and response bodies must be in JSON format.


## Installation

Obtain the following libraries:
```
go get github.com/lib/pq
go get github.com/gorilla/mux
go get github.com/dgrijalva/jwt-go
```

Optional. Define authentication secret in the environment:

```
export AUTH_SECRET="big_secret"
go build
./server.exe
```
Test plain status endpoint:

```
curl -I localhost:8080/api/v1/status
```

Obtain authentication token:

```
curl localhost:8080/api/v1/auth -X POST -H "Content-Type: application/json" --data '{"username":"xyz","password":"xyz"}'
```

Copy token and use it for future requests in the header. For example:

```
curl localhost:8080/api/v1/patients \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODU3ODIwMjUsInVzZXJuYW1lIjoieHl6In0.vOZ9iBqn325Sn8_EnDod2emqPnGxTdssD17qlGP4xEg"
```


