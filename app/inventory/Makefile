mysqlUser=inventory
mysqlPassword=inventory
mysqlDatabase=inventory
mysqlHost=localhost
mysqlPort=3308

mock-gen: goa-gen-vehicle
	go generate ./...

# Generates all the mocks for this Inventory Context including the generation of common pkg mocks
mock-gen-with-pkg:
	$(MAKE) mock-gen
	$(MAKE) -C ../../pkg mock-gen

mock-clean:
	@find . -type f -name '*_mock.go' -delete
	@cd ../../pkg && $(MAKE) mock-clean

unit-test: goa-gen-vehicle mock-gen
	@CGO_ENABLED=1 go test ./... -tags=unit -race

integration-test: export MYSQL_USER=$(mysqlUser)
integration-test: export MYSQL_PASSWORD=$(mysqlPassword)
integration-test: export MYSQL_DATABASE=$(mysqlDatabase)
integration-test: export MYSQL_HOST=$(mysqlHost)
integration-test: export MYSQL_PORT=$(mysqlPort)
integration-test: goa-gen-vehicle mock-gen
	@CGO_ENABLED=1 go test ./... -tags=integration -p 1 -race

# OPENAPI 3 Generation approaches:

goa-gen-vehicle:
	goa gen vehicle-sharing-go/app/inventory/internal/vehicle/handler/goa/design \
		-o internal/vehicle/handler/goa

swaggo-gen-vehicle:
	swag init \
		-g api_docs.go \
	  	-d internal/vehicle/handler/swago-annotations,internal/vehicle/projection,internal/vehicle/command,$(PKG_DIR) \
	  	-o internal/vehicle/handler/swago-annotations/docs \
	  	--parseDependency \
	  	--v3.1 true

run-rest-api-goa:
	go run ./cmd/vehicle/rest-api --server-library=goa

run-rest-api-net-http:
	go run ./cmd/vehicle/rest-api --server-library=net-http

run-rest-api-gin:
	go run ./cmd/vehicle/rest-api --server-library=gin