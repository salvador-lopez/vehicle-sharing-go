mock-gen:
	go generate ./...

mock-clean:
	@find . -type f -name '*_mock.go' -delete

unit-test: mock-gen
	@CGO_ENABLED=1 go test -race -covermode atomic -short ./...