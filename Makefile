mock-gen:
	go generate ./...

mock-clean:
	@find . -type f -name '*_mock.go' -delete

unit-test-inventory: mock-gen
	@CGO_ENABLED=1 go test ./... -tags=unit_inventory -race -covermode atomic

integration-test-inventory: export MYSQL_USER=inventory
integration-test-inventory: export MYSQL_PASSWORD=inventory
integration-test-inventory: export MYSQL_DATABASE=inventory
integration-test-inventory: export MYSQL_HOST=localhost
integration-test-inventory: export MYSQL_PORT=3308

integration-test-inventory: mock-gen
	CGO_ENABLED=1 go test ./... -tags=integration_inventory -p 1 -race

goa-gen-inventory-vehicle:
	goa gen vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/design -o internal/inventory/vehicle/infrastructure/controller