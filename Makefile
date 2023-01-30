mock-gen:
	go generate ./...

mock-clean:
	@find . -type f -name '*_mock.go' -delete

unit-test: mock-gen
	@CGO_ENABLED=1 go test ./... -tags=unit -race -covermode atomic -short

integration-test: export MYSQL_USER=user
integration-test: export MYSQL_PASSWORD=pass
integration-test: export MYSQL_HOST=mysql
integration-test: export MYSQL_DATABASE=vehicle-sharing
integration-test: export MYSQL_HOST=localhost
integration-test: export MYSQL_PORT=3306

integration-test: mock-gen
	CGO_ENABLED=1 go test ./... -tags=integration -p 1 -race