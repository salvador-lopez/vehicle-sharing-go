// Code generated by goa v3.11.1, DO NOT EDIT.
//
// car service
//
// Command:
// $ goa gen
// vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/design
// -o internal/inventory/vehicle/infrastructure/controller

package car

import (
	"context"
	carviews "vehicle-sharing-go/internal/inventory/vehicle/infrastructure/handler/gen/car/views"

	goa "goa.design/goa/v3/pkg"
)

// The car service performs operations on car vehicles inventory
type Service interface {
	// Create implements create.
	Create(context.Context, *CreatePayload) (err error)
	// Get implements get.
	Get(context.Context, *GetPayload) (res *CarResource, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "car"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"create", "get"}

// CarResource is the result type of the car service get method.
type CarResource struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	Color     string
	VinData   *VinData
}

// CreatePayload is the payload type of the car service create method.
type CreatePayload struct {
	// Car id in uuid format
	ID  string
	Vin Vin
	// Car color
	Color string
}

// GetPayload is the payload type of the car service get method.
type GetPayload struct {
	// Car id in uuid format
	ID string
}

type Vin string

// decoded vin data
type VinData struct {
	Vin           Vin
	Country       *string
	Manufacturer  *string
	Brand         *string
	EngineSize    *string
	FuelType      *string
	Model         *string
	Year          *string
	AssemblyPlant *string
	SN            *string
}

// MakeInternal builds a goa.ServiceError from an error.
func MakeInternal(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "internal", false, false, false)
}

// MakeConflict builds a goa.ServiceError from an error.
func MakeConflict(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "conflict", false, false, false)
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "notFound", false, false, false)
}

// NewCarResource initializes result type CarResource from viewed result type
// CarResource.
func NewCarResource(vres *carviews.CarResource) *CarResource {
	return newCarResource(vres.Projected)
}

// NewViewedCarResource initializes viewed result type CarResource from result
// type CarResource using the given view.
func NewViewedCarResource(res *CarResource, view string) *carviews.CarResource {
	p := newCarResourceView(res)
	return &carviews.CarResource{Projected: p, View: "default"}
}

// newCarResource converts projected type CarResource to service type
// CarResource.
func newCarResource(vres *carviews.CarResourceView) *CarResource {
	res := &CarResource{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.CreatedAt != nil {
		res.CreatedAt = *vres.CreatedAt
	}
	if vres.UpdatedAt != nil {
		res.UpdatedAt = *vres.UpdatedAt
	}
	if vres.Color != nil {
		res.Color = *vres.Color
	}
	if vres.VinData != nil {
		res.VinData = newVinData(vres.VinData)
	}
	return res
}

// newCarResourceView projects result type CarResource to projected type
// CarResourceView using the "default" view.
func newCarResourceView(res *CarResource) *carviews.CarResourceView {
	vres := &carviews.CarResourceView{
		ID:        &res.ID,
		CreatedAt: &res.CreatedAt,
		UpdatedAt: &res.UpdatedAt,
		Color:     &res.Color,
	}
	if res.VinData != nil {
		vres.VinData = newVinDataView(res.VinData)
	}
	return vres
}

// newVinData converts projected type VinData to service type VinData.
func newVinData(vres *carviews.VinDataView) *VinData {
	res := &VinData{
		Country:       vres.Country,
		Manufacturer:  vres.Manufacturer,
		Brand:         vres.Brand,
		EngineSize:    vres.EngineSize,
		FuelType:      vres.FuelType,
		Model:         vres.Model,
		Year:          vres.Year,
		AssemblyPlant: vres.AssemblyPlant,
		SN:            vres.SN,
	}
	if vres.Vin != nil {
		res.Vin = Vin(*vres.Vin)
	}
	return res
}

// newVinDataView projects result type VinData to projected type VinDataView
// using the "default" view.
func newVinDataView(res *VinData) *carviews.VinDataView {
	vres := &carviews.VinDataView{
		Country:       res.Country,
		Manufacturer:  res.Manufacturer,
		Brand:         res.Brand,
		EngineSize:    res.EngineSize,
		FuelType:      res.FuelType,
		Model:         res.Model,
		Year:          res.Year,
		AssemblyPlant: res.AssemblyPlant,
		SN:            res.SN,
	}
	vin := carviews.VinView(res.Vin)
	vres.Vin = &vin
	return vres
}
