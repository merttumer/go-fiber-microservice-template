package httptansport

import (
	"fmt"
	"reflect"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/gofiber/fiber/v2"
	dummyservice "github.com/merttumer/go-fiber-microservice-template"
	"github.com/merttumer/go-fiber-microservice-template/pkg/apierror"
	"github.com/merttumer/go-fiber-microservice-template/pkg/endpoints"
)

type DecodeRequestFunc func(*fiber.Ctx) (interface{}, error)
type EncodeResponseFunc func(*fiber.Ctx, interface{}) error

func MakeHttpHandler(l log.Logger, hs dummyservice.Service) *fiber.App {
	es := endpoints.MakeEndpoints(hs)
	app := fiber.New(fiber.Config{})

	app.Route("/sum", sumRouter(es, l))

	return app
}

func sumRouter(es endpoints.Endpoints, l log.Logger) func(r fiber.Router) {
	return func(r fiber.Router) {
		r.Add("GET", "/", makeSumHandler(es.SumEndpoint))
	}
}

func makeSumHandler(e endpoint.Endpoint) fiber.Handler {
	return newDefaultHandler(e, dummyservice.SumRequest{})
}

func newDefaultHandler(e endpoint.Endpoint, emptyReq interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := makeDecoder(emptyReq)(c)
		if err != nil {
			return err
		}

		res, err := e(c.Context(), req)

		if err != nil {
			errorEncoder(c, err)
			return err
		}

		return makeEncoder()(c, res)
	}
}

func makeDecoder(emptyReq interface{}) DecodeRequestFunc {
	return func(c *fiber.Ctx) (interface{}, error) {
		req := reflect.New(reflect.TypeOf(emptyReq)).Interface()

		if err := c.ParamsParser(req); err != nil {
			return nil, fmt.Errorf("error parsing params: %w", err)
		}

		if err := c.ReqHeaderParser(req); err != nil {
			return nil, fmt.Errorf("error parsing header: %w", err)
		}

		if err := c.QueryParser(req); err != nil {
			return nil, fmt.Errorf("error parsing query: %w", err)
		}

		if len(c.Body()) > 0 {
			if err := c.BodyParser(req); err != nil {
				return nil, fmt.Errorf("error parsing body: %w", err)
			}
		}

		return req, nil
	}
}

func makeEncoder() EncodeResponseFunc {
	return func(c *fiber.Ctx, response interface{}) error {
		res, ok := response.(dummyservice.Response)
		if !ok {
			return fmt.Errorf("returned response does not implement Response interface")
		}

		if err := res.APIError(); err != nil {
			return err
		}

		c.JSON(response)
		return nil
	}
}

type ErrorResponse struct {
	Data   interface{} `json:"data"`
	Result error       `json:"result"`
}

func errorEncoder(c *fiber.Ctx, err error) error {
	var apiErr *apierror.APIError

	switch v := err.(type) {
	case *apierror.APIError:
		apiErr = v
		c.Status(apiErr.StatusCode)
	default:
		apiErr = apierror.NewInternalServerError()
		c.Status(apiErr.StatusCode)
	}
	c.JSON(ErrorResponse{
		Data:   nil,
		Result: apiErr,
	})
	return nil
}
