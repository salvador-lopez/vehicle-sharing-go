// Code generated by goa v3.11.1, DO NOT EDIT.
//
// car views
//
// Command:
// $ goa gen
// vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/design
// -o internal/inventory/vehicle/infrastructure/controller

package views

import (
	goa "goa.design/goa/v3/pkg"
)

// CarResource is the viewed result type that is projected based on a view.
type CarResource struct {
	// Type to project
	Projected *CarResourceView
	// View to render
	View string
}

// CarResourceView is a type that runs validations on a projected type.
type CarResourceView struct {
	ID        *string
	CreatedAt *string
	UpdatedAt *string
	Color     *string
	VinData   *VinDataView
}

// VinDataView is a type that runs validations on a projected type.
type VinDataView struct {
	Vin           *VinView
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

// VinView is a type that runs validations on a projected type.
type VinView string

var (
	// CarResourceMap is a map indexing the attribute names of CarResource by view
	// name.
	CarResourceMap = map[string][]string{
		"default": {
			"id",
			"createdAt",
			"updatedAt",
			"color",
			"vinData",
		},
	}
	// VinDataMap is a map indexing the attribute names of VinData by view name.
	VinDataMap = map[string][]string{
		"default": {
			"vin",
			"country",
			"manufacturer",
			"brand",
			"engineSize",
			"fuelType",
			"model",
			"year",
			"assemblyPlant",
			"SN",
		},
	}
)

// ValidateCarResource runs the validations defined on the viewed result type
// CarResource.
func ValidateCarResource(result *CarResource) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateCarResourceView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateCarResourceView runs the validations defined on CarResourceView
// using the "default" view.
func ValidateCarResourceView(result *CarResourceView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.CreatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("createdAt", "result"))
	}
	if result.UpdatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("updatedAt", "result"))
	}
	if result.Color == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("color", "result"))
	}
	if result.ID != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.id", *result.ID, goa.FormatUUID))
	}
	if result.CreatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.createdAt", *result.CreatedAt, goa.FormatDateTime))
	}
	if result.UpdatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.updatedAt", *result.UpdatedAt, goa.FormatDateTime))
	}
	if result.VinData != nil {
		if err2 := ValidateVinDataView(result.VinData); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateVinDataView runs the validations defined on VinDataView using the
// "default" view.
func ValidateVinDataView(result *VinDataView) (err error) {
	if result.Vin == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("vin", "result"))
	}
	return
}

// ValidateVinView runs the validations defined on VinView.
func ValidateVinView(result VinView) (err error) {

	return
}
