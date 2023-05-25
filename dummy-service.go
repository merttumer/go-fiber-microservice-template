package dummyservice

import (
	"context"

	"github.com/merttumer/go-fiber-microservice-template/pkg/apierror"
)

type Service interface {
	Sum(context.Context, SumRequest) SumResponse
}

type Request interface {
}

type Response interface {
	Localize() interface{}
	APIError() error
}

// compile time proof of request and response implementations
var (
	_ Request  = (*SumRequest)(nil)
	_ Response = (*SumResponse)(nil)
)

type (
	SumRequest struct {
		A int `query:"a" json:"-"`
		B int `query:"b" json:"-"`
	}

	SumResponse struct {
		Result *apierror.APIError `json:"-"`
		Res    string             `json:"res"`
	}
)

func (s SumResponse) Localize() interface{} {
	// localization logic
	return s
}

func (s SumResponse) APIError() error {
	if s.Result == nil {
		return nil
	}
	return s.Result
}
