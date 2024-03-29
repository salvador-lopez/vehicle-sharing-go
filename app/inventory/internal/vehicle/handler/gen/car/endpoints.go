// Code generated by goa v3.11.1, DO NOT EDIT.
//
// car endpoints
//
// Command:
// $ goa gen
// vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/design
// -o internal/inventory/vehicle/infrastructure/controller

package car

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "car" service endpoints.
type Endpoints struct {
	Create goa.Endpoint
	Get    goa.Endpoint
}

// NewEndpoints wraps the methods of the "car" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Create: NewCreateEndpoint(s),
		Get:    NewGetEndpoint(s),
	}
}

// Use applies the given middleware to all the "car" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Create = m(e.Create)
	e.Get = m(e.Get)
}

// NewCreateEndpoint returns an endpoint function that calls the method
// "create" of service "car".
func NewCreateEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*CreatePayload)
		return nil, s.Create(ctx, p)
	}
}

// NewGetEndpoint returns an endpoint function that calls the method "get" of
// service "car".
func NewGetEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*GetPayload)
		res, err := s.Get(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedCarResource(res, "default")
		return vres, nil
	}
}
