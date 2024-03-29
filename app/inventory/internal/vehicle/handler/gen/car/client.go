// Code generated by goa v3.11.1, DO NOT EDIT.
//
// car client
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

// Client is the "car" service client.
type Client struct {
	CreateEndpoint goa.Endpoint
	GetEndpoint    goa.Endpoint
}

// NewClient initializes a "car" service client given the endpoints.
func NewClient(create, get goa.Endpoint) *Client {
	return &Client{
		CreateEndpoint: create,
		GetEndpoint:    get,
	}
}

// Create calls the "create" endpoint of the "car" service.
// Create may return the following errors:
//   - "conflict" (type *goa.ServiceError)
//   - "internal" (type *goa.ServiceError)
//   - error: internal error
func (c *Client) Create(ctx context.Context, p *CreatePayload) (err error) {
	_, err = c.CreateEndpoint(ctx, p)
	return
}

// Get calls the "get" endpoint of the "car" service.
// Get may return the following errors:
//   - "notFound" (type *goa.ServiceError)
//   - "internal" (type *goa.ServiceError)
//   - error: internal error
func (c *Client) Get(ctx context.Context, p *GetPayload) (res *CarResource, err error) {
	var ires interface{}
	ires, err = c.GetEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*CarResource), nil
}
