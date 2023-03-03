package design

import . "goa.design/goa/v3/dsl"

var vin = Type("vin", String)

var vinData = ResultType("vin_data", func() {
	Description("decoded vin data")
	TypeName("vinData")
	Attributes(func() {
		Attribute("vin", vin)
		Attribute("country", String)
		Attribute("manufacturer", String)
		Attribute("brand", String)
		Attribute("engineSize", String)
		Attribute("fuelType", String)
		Attribute("model", String)
		Attribute("year", String)
		Attribute("assemblyPlant", String)
		Attribute("SN", String)
		Required("vin")
	})
})

var carResource = ResultType("inventory.vehicle.car", func() {
	Description("Car resource containing vin decoded data")
	TypeName("carResource")
	Attributes(func() {
		Attribute("id", String, func() {
			Format(FormatUUID)
		})
		Attribute("createdAt", String, func() {
			Format(FormatDateTime)
		})
		Attribute("updatedAt", String, func() {
			Format(FormatDateTime)
		})
		Attribute("color", String)
		Attribute("vinData", vinData)
		Required("id", "createdAt", "updatedAt", "color", "vinData")
	})
})

// Service describes a service
var _ = Service("car", func() {
	HTTP(func() {
		Path("/cars")
		Response("internal", StatusInternalServerError)
	})

	Error("internal")

	Description("The car service performs operations on car vehicles inventory")
	Method("create", func() {
		Payload(func() {
			Attribute("id", String, "Car id in uuid format", func() {
				Format(FormatUUID)
			})
			Attribute("vin", vin)
			Attribute("color", String, "Car color")
			Required("id", "vin", "color")
		})
		HTTP(func() {
			POST("/")
			Response(StatusCreated)
			Response(StatusBadRequest)
			Response("conflict", StatusConflict)
		})

		Error("conflict")
	})

	Method("get", func() {
		Payload(func() {
			Attribute("id", String, "Car id in uuid format", func() {
				Format(FormatUUID)
			})
			Required("id")
		})
		HTTP(func() {
			GET("/{id}")
			Response(StatusOK)
			Response(StatusBadRequest)
			Response("notFound", StatusNotFound)
		})

		Error("notFound")

		Result(carResource)
	})
})
