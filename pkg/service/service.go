package service

import (
	"context"
	"fmt"

	dummyservice "github.com/merttumer/go-fiber-microservice-template"
)

type service struct{}

func NewService() dummyservice.Service {
	return &service{}
}

// Sum implements Service
func (s *service) Sum(_ context.Context, req dummyservice.SumRequest) dummyservice.SumResponse {
	return dummyservice.SumResponse{
		Res: fmt.Sprintf("%d + %d = %d", req.A, req.B, req.A+req.B),
	}
}
