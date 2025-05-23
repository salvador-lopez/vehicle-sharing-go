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
	"vehicle-sharing-go/app/inventory/cmd/vehicle/rest-api/gin"
	"vehicle-sharing-go/app/inventory/cmd/vehicle/rest-api/goa"
	nethttp "vehicle-sharing-go/app/inventory/cmd/vehicle/rest-api/net-http"
	"vehicle-sharing-go/app/inventory/internal/vehicle/database/gorm/model"
	modelpkg "vehicle-sharing-go/pkg/database/gorm/model"

	"github.com/google/uuid"

	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/database/gorm"
	gormpkg "vehicle-sharing-go/pkg/database/gorm"
	commandpkg "vehicle-sharing-go/pkg/domain/event"
)

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		serverLibrary = flag.String("server-library", "", "define which server library to use")
		hostF         = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF       = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF     = flag.String("http-port", "", "HTTP port (overrides host HTTP port specified in service design)")

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

	// Setup logger.
	logger := log.New(os.Stderr, fmt.Sprintf("[inventory-vehicles-rest-api-%s] ", *serverLibrary), log.Ltime)

	dbConn, err := gormpkg.NewConnectionFromConfig(&gormpkg.Config{
		UserName:     *dbUser,
		Password:     *dbPwd,
		DatabaseName: *dbName,
		Host:         *dbHost,
		Port:         *dbPort,
		Logger:       logger,
		LogQueries:   *dbDebug,
	})
	if err != nil {
		logger.Fatalf("failed to create db connection: %v", err)
	}

	// Initialize Write Repositories
	carRepo := gorm.NewCarRepository(dbConn)
	err = dbConn.Db().AutoMigrate(&model.Car{})
	if err != nil {
		logger.Fatalf("AutoMigrate Car model failed: %v", err)
	}

	// Initialize AggregateRoot Domain Events Publisher
	outboxRepo := gormpkg.NewOutboxRepository(dbConn)
	err = dbConn.Db().AutoMigrate(&modelpkg.OutboxRecord{})
	if err != nil {
		logger.Fatalf("AutoMigrate Outbox model failed: %v", err)
	}
	aggRootEventPublisher := commandpkg.NewAgRootEventPublisher(outboxRepo)

	// Initialize Query Services
	carQueryService := gorm.NewCarProjectionRepository(dbConn.Db())

	// Initialize commandHandlers
	idGenFn := uuid.New
	nowFn := time.Now

	createCarHandler := command.NewCreateCarHandler(idGenFn, nowFn, carRepo, dbConn, aggRootEventPublisher)

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

			shutdownHook := registerShutdownHook(ctx, logger)

			switch *serverLibrary {
			case "net-http":
				nethttp.HandleHTTPServer(ctx, shutdownHook, u, carQueryService, createCarHandler, &wg, errc, logger)
			case "goa":
				goa.HandleHTTPServer(ctx, u, carQueryService, createCarHandler, &wg, errc, logger, *dbgF)
			case "gin":
				gin.HandleHTTPServer(shutdownHook, u, carQueryService, createCarHandler, &wg, errc, logger, *dbgF)
			default:
				logger.Println("No server library defined, defaulting to gin")
				gin.HandleHTTPServer(shutdownHook, u, carQueryService, createCarHandler, &wg, errc, logger, *dbgF)
			}
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
