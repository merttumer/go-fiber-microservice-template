package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/log"
	dummyservice "github.com/merttumer/go-fiber-microservice-template"
	"github.com/merttumer/go-fiber-microservice-template/pkg/middleware/metrics"
	"github.com/merttumer/go-fiber-microservice-template/pkg/service"
	httptransport "github.com/merttumer/go-fiber-microservice-template/pkg/transport/http"
)

type CustomLoggerAdapter struct {
	logger log.Logger
}

func (c *CustomLoggerAdapter) Printf(format string, args ...interface{}) {
	(c.logger).Log("msg", fmt.Sprintf(format, args...))
}

func main() {
	var l log.Logger
	{
		l = log.NewLogfmtLogger(os.Stdout)
		l = log.With(l, "time", log.DefaultTimestampUTC)
	}

	var svc dummyservice.Service
	{
		svc = service.NewService()
	}

	var mm dummyservice.Service
	{
		mm = metrics.NewMetricsMiddleware()(svc)
	}

	app := httptransport.MakeHttpHandler(l, mm)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			app.Shutdown()
		}
	}()

	// Wait for a signal to shut down the server
	<-sig

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shut down the server gracefully
	if err := app.ShutdownWithContext(ctx); err != nil {
		fmt.Printf("%s", err.Error())
	}
}
