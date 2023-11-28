package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"

	"vehicle-sharing-go/internal/inventory/vehicle/application/command"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/gen/car"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/rest"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm"
	inmemory "vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/in-memory"
)

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF = flag.String("http-port", "", "HTTP port (overrides host HTTP port specified in service design)")

		dbUser  = flag.String("db-user", "inventory", "database user")
		dbPwd   = flag.String("db-password", "inventory", "database password")
		dbName  = flag.String("db-name", "inventory", "database name")
		dbHost  = flag.String("db-host", "localhost", "database host")
		dbPort  = flag.Int("db-port", 3308, "database port")
		dbDebug = flag.Bool("db-debug", false, "database debug mode")

		secureF = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF    = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	// Setup logger. Replace logger with your own log package of choice.
	var (
		logger *log.Logger
	)
	{
		logger = log.New(os.Stderr, "[inventoryvehicles] ", log.Ltime)
	}

	cfg := &gorm.Config{
		UserName:     *dbUser,
		Password:     *dbPwd,
		DatabaseName: *dbName,
		Host:         *dbHost,
		Port:         *dbPort,
		Logger:       logger,
		LogQueries:   *dbDebug,
	}

	dbConn, err := gorm.NewConnectionFromConfig(cfg)
	if err != nil {
		logger.Fatalf("failed to create db connection: %v", err)
	}

	// Initialize the services.
	var (
		carSvc car.Service
	)
	{
		carSvc = rest.NewCarController(
			command.NewCreateCarHandler(uuid.New, time.Now, dbConn, dbConn),
			inmemory.NewCarQueryService(),
		)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		carEndpoints *car.Endpoints
	)
	{
		carEndpoints = car.NewEndpoints(carSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "localhost":
		{
			addr := "http://localhost:8088/api/inventory/vehicles"
			u, err := url.Parse(addr)
			if err != nil {
				logger.Fatalf("invalid URL %#v: %s\n", addr, err)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h, _, err := net.SplitHostPort(u.Host)
				if err != nil {
					logger.Fatalf("invalid URL %#v: %s\n", u.Host, err)
				}
				u.Host = net.JoinHostPort(h, *httpPortF)
			} else if u.Port() == "" {
				u.Host = net.JoinHostPort(u.Host, "80")
			}
			handleHTTPServer(ctx, u, carEndpoints, &wg, errc, logger, *dbgF)
		}

	default:
		logger.Fatalf("invalid host argument: %q (valid hosts: localhost)\n", *hostF)
	}

	// Wait for signal.
	logger.Printf("exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Println("exited")
}
