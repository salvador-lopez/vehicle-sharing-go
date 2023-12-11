mysqlUser=inventory
mysqlPassword=inventory
mysqlDatabase=inventory
mysqlHost=localhost
mysqlPort=3308

mock-gen:
	go generate ./...

mock-clean:
	@find . -type f -name '*_mock.go' -delete

unit-test: mock-gen
	@CGO_ENABLED=1 go test ./... -tags=unit -race

unit-test-inventory: mock-gen
	@CGO_ENABLED=1 go test ./... -tags=unit_inventory -race

integration-test: export MYSQL_USER=$(mysqlUser)
integration-test: export MYSQL_PASSWORD=$(mysqlPassword)
integration-test: export MYSQL_DATABASE=$(mysqlDatabase)
integration-test: export MYSQL_HOST=$(mysqlHost)
integration-test: export MYSQL_PORT=$(mysqlPort)
integration-test: mock-gen
	@CGO_ENABLED=1 go test ./... -tags=integration -race

integration-test-inventory: export MYSQL_USER=$(mysqlUser)
integration-test-inventory: export MYSQL_PASSWORD=$(mysqlPassword)
integration-test-inventory: export MYSQL_DATABASE=$(mysqlDatabase)
integration-test-inventory: export MYSQL_HOST=$(mysqlHost)
integration-test-inventory: export MYSQL_PORT=$(mysqlPort)
integration-test-inventory: mock-gen
	CGO_ENABLED=1 go test ./... -tags=integration_inventory -p 1 -race

goa-gen-inventory-vehicle:
	goa gen vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/design -o internal/inventory/vehicle/infrastructure/controller