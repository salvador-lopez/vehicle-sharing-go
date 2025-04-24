# vehicle-sharing-go
This is an example of a vehicle sharing system using the monorepo multiapplication approach in golang

## 🚀 Run the bash commands in your local terminal if you want to make them permanent. i.e. when the commands are installing a library

## DEV environment
### Pre-requisites
Install docker & docker compose in your host SO:
- [Get Docker](https://docs.docker.com/get-started/get-docker/)
- [Install Docker Compose](https://docs.docker.com/compose/install/)

To generate OPENAPI 3 DOCS:

- When using [GOA Design](https://goa.design/) as an HTTP Server
```bash
go install goa.design/goa/v3/cmd/goa@latest
export PATH="$PATH:$(go env GOPATH)/bin"
source ~/.zshrc  # if you're using ohmyzsh or source ~/.bashrc if you are using regular bash
```
To check that goa is properly installed:
```bash
goa version
```
And then to generate the openapi doc files and the needed generated resources from project root path:
```bash
make goa-gen
```
- When using [net/http go stlib](https://pkg.go.dev/net/http) as an HTTP Server
```bash
go install github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc4 # We need the v2 version as the lastest stable 1.6.4 not generate openapi3 spec
export PATH="$PATH:$(go env GOPATH)/bin"
source ~/.zshrc  # if you're using ohmyzsh or source ~/.bashrc if you are using regular bash
```
To check that swag is properly installed:
```bash
swag --version
```
And then to generate the openapi doc and generated files from project root path:
```bash
make swaggo-gen
```

> 💡 **TIP:** You can also create/update your goa/swag files for a specific Context, i.e. `inventory`. From project root path:
```bash
cd app/inventory
make goa-gen-vehicle 
make swaggo-gen-vehicle

```
To create and/or clean mocks we need to install mockgen:
```bash
go install github.com/golang/mock/mockgen@v1.6.0
export PATH="$PATH:$(go env GOPATH)/bin"
source ~/.zshrc  # if you're using ohmyzsh or source ~/.bashrc if you are using regular bash
```

### Setting up DEV environment
From project root path:
```bash
docker compose up -d
```
### Running DEV environment
We can run the different services needed to run our application in localhost.
If you code with Goland You can use the [Project Running Configurations](https://www.jetbrains.com/help/go/run-debug-configuration.html):
- REST API Server:
  - [With GOA Design](.run/inventory/vehicles/rest-api-goa.run.xml)
  - [With net/http and swaggo](.run/inventory/vehicles/rest-api-nethttp.run.xml)
- Domain Event Consumer:
  - [With Kafka](.run/inventory/vehicles/domain-event-consumer-kafka.run.xml)
- Message Relay:
  - [With Mysql Binary Log](.run/inventory/vehicles/message-relay-mysql-binlog.run.xml)


### Testing
The project contains two different kind of automated tests, unit and integration. We can run the tests for all the Contexts
Using the Makefile from the project root path:
```bash
make unit-test
```
```bash
make integration-test
```

You can also run tests for a specific Context:
- Inventory:
```bash
cd app/inventory
make unit-test
make integration-test
```