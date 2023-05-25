package metrics

import (
	"context"
	"log"
	"time"

	dummyservice "github.com/merttumer/go-fiber-microservice-template"
	"github.com/merttumer/go-fiber-microservice-template/pkg/middleware"
)

type MetricsMiddleware struct {
	next dummyservice.Service
}

// MetricsMiddleware is a middleware that collects metrics
func NewMetricsMiddleware() middleware.Middleware {
	return func(next dummyservice.Service) dummyservice.Service {
		return &MetricsMiddleware{
			next: next,
		}
	}
}

// Sum implements dummyservice.Service
func (mw *MetricsMiddleware) Sum(ctx context.Context, svc dummyservice.SumRequest) dummyservice.SumResponse {
	// log the time it took to execute the method
	startTime := time.Now()

	res := mw.next.Sum(ctx, svc)

	duration := time.Since(startTime)
	log.Printf("Time taken to execute Sum method: %v", duration)

	return res
}
