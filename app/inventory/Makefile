mysqlUser=inventory
mysqlPassword=inventory
mysqlDatabase=inventory
mysqlHost=localhost
mysqlPort=3308

mock-gen:
	go generate ./...

# Generates all the mocks for this Inventory Context including the generation of common pkg mocks
mock-gen-with-pkg:
	$(MAKE) mock-gen
	$(MAKE) -C ../../pkg mock-gen

mock-clean:
	@find . -type f -name '*_mock.go' -delete
	@cd ../../pkg && $(MAKE) mock-clean

unit-test:
	@CGO_ENABLED=1 go test ./... -tags=unit -race

integration-test: export MYSQL_USER=$(mysqlUser)
integration-test: export MYSQL_PASSWORD=$(mysqlPassword)
integration-test: export MYSQL_DATABASE=$(mysqlDatabase)
integration-test: export MYSQL_HOST=$(mysqlHost)
integration-test: export MYSQL_PORT=$(mysqlPort)
integration-test: mock-gen
	@CGO_ENABLED=1 go test ./... -tags=integration -p 1 -race

# OPENAPI 3 Generation approaches:

goa-gen-vehicle:
	goa gen vehicle-sharing-go/app/inventory/internal/vehicle/handler/goa/design \
		-o internal/vehicle/handler/goa

swaggo-gen-vehicle:
	swag init \
	  -g http.go \
	  -d cmd/vehicle/rest-api/net-http,internal/vehicle/handler/net-http/rest \
	  -o internal/vehicle/handler/net-http/docs \
	  --parseDependency \
	  --v3.1 true