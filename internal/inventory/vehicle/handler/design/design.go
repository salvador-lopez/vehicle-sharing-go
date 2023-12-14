package design

import . "goa.design/goa/v3/dsl"

// API describes the global properties of the API server.
var _ = API("inventory vehicles", func() {
	Title("Inventory Service")
	Description("HTTP service to interact with inventory bounded context")
	Server("inventory/vehicles", func() {
		Host("localhost", func() { URI("http://localhost:8088/api/inventory/vehicles") })
	})
	HTTP(func() {
		Path("api/inventory/vehicles")
	})
})
