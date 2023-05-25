package middleware

import dummyservice "github.com/merttumer/go-fiber-microservice-template"

type Middleware func(dummyservice.Service) dummyservice.Service
