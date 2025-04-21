mysqlUser=inventory
mysqlPassword=inventory
mysqlDatabase=inventory
mysqlHost=localhost
mysqlPort=3308

mock-gen:
	go generate ./...

mock-clean:
	@find . -type f -name '*_mock.go' -delete

unit-test:
	@CGO_ENABLED=1 go test ./... -tags=unit -race

integration-test: export MYSQL_USER=$(mysqlUser)
integration-test: export MYSQL_PASSWORD=$(mysqlPassword)
integration-test: export MYSQL_DATABASE=$(mysqlDatabase)
integration-test: export MYSQL_HOST=$(mysqlHost)
integration-test: export MYSQL_PORT=$(mysqlPort)
integration-test: mock-gen
	@CGO_ENABLED=1 go test ./... -tags=integration -p 1 -race