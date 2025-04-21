package design

import . "goa.design/goa/v3/dsl"

// API describes the global properties of the API server.
var _ = API("inventory vehicles", func() {
	Title("Inventory Vehicles Service")
	Description("HTTP service to interact with inventory vehicles bounded context")
	Server("inventory/vehicles", func() {
		Host("localhost", func() { URI("http://localhost:8088") })
	})
	HTTP(func() {
		Path("api/inventory/vehicles")
	})
})
