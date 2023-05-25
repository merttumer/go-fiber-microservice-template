package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	dummyservice "github.com/merttumer/go-fiber-microservice-template"
)

type Endpoints struct {
	SumEndpoint endpoint.Endpoint
}

func makeSumEndpoint(svc dummyservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dummyservice.SumRequest)

		res := svc.Sum(ctx, *req)

		return res, nil
	}
}

func MakeEndpoints(svc dummyservice.Service) Endpoints {
	return Endpoints{
		SumEndpoint: makeSumEndpoint(svc),
	}
}
