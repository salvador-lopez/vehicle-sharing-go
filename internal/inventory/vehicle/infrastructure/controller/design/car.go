package design

import . "goa.design/goa/v3/dsl"

// Service describes a service
var _ = Service("car", func() {
	Error("badRequest")
	Error("conflict")
	Error("internal")

	HTTP(func() {
		Path("/cars")
		Response("badRequest", StatusBadRequest)
		Response("conflict", StatusConflict)
		Response("internal", StatusInternalServerError)
	})

	Description("The car service performs operations on car vehicles inventory")
	Method("create", func() {
		Payload(func() {
			Attribute("id", String, "Car id in uuid format", func() {
				Format(FormatUUID)
			})
			Attribute("vin", String, "Vehicle Identification Number")
			Attribute("color", String, "Car color")
			Required("id", "vin", "color")
		})
		HTTP(func() {
			POST("/")
			Response(StatusCreated)
		})
	})
})
