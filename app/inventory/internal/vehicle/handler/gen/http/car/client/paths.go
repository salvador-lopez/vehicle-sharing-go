// Code generated by goa v3.11.1, DO NOT EDIT.
//
// HTTP request path constructors for the car service.
//
// Command:
// $ goa gen
// vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/design
// -o internal/inventory/vehicle/infrastructure/controller

package client

import (
	"fmt"
)

// CreateCarPath returns the URL path to the car service create HTTP endpoint.
func CreateCarPath() string {
	return "/api/inventory/vehicles/cars"
}

// GetCarPath returns the URL path to the car service get HTTP endpoint.
func GetCarPath(id string) string {
	return fmt.Sprintf("/api/inventory/vehicles/cars/%v", id)
}
