package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"vehicle-sharing-go/internal/inventory/vehicle/database/gorm"
	"vehicle-sharing-go/internal/inventory/vehicle/domain/event"
	"vehicle-sharing-go/internal/inventory/vehicle/projection"
	gorm2 "vehicle-sharing-go/pkg/database/gorm"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the consumer",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		c, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":     "localhost:19092",
			"broker.address.family": "v4",
			"group.id":              "inventory-vehicles",
			"auto.offset.reset":     "earliest",
		})

		// Setup logger. Replace logger with your own log package of choice.
		logger := log.New(os.Stderr, "[inventory-vehicles-domain-event-consumer] ", log.Ltime)

		dbConn, err := gorm2.NewConnectionFromConfig(&gorm2.Config{
			UserName:     "inventory",
			Password:     "inventory",
			DatabaseName: "inventory",
			Host:         "localhost",
			Port:         3308,
			Logger:       logger,
			LogQueries:   false,
		})
		if err != nil {
			logger.Fatalf("failed to create db connection: %v", err)
		}

		// Initialize Write Repositories
		carRepo := gorm.NewCarProjectionRepository(dbConn.Db())

		carProjector := projection.NewCarProjector(vinDecoderFake{}, carRepo)

		if err != nil {
			panic(err)
		}

		defer c.Close()

		err = c.SubscribeTopics([]string{"inventory-vehicles-car"}, nil)
		if err != nil {
			panic(err)
		}

		go func() {
			for {
				msg, err := c.ReadMessage(time.Second)
				if err == nil {
					logger.Printf("Message on %s with AggregateID %s: %s\n", msg.TopicPartition, string(msg.Key), string(msg.Value))

					aggregateID, err := uuid.Parse(string(msg.Key))
					if err != nil {
						logger.Printf("Error Unmarshalling message aggregateID: %v (%v)\n", err, msg)
						errc <- err
					}

					var payload *event.CarCreatedPayload
					err = json.Unmarshal(msg.Value, &payload)
					if err != nil {
						logger.Printf("Error Unmarshalling message: %v (%v)\n", err, msg)
						errc <- err
					}

					err = carProjector.ProjectCarCreated(ctx, aggregateID, payload)
					if err != nil {
						logger.Printf("Error Projecting Event %s: %v (%v)\n", aggregateID, err, payload)
						errc <- err
					}
				} else if !err.(kafka.Error).IsTimeout() {
					// The client will automatically try to recover from all errors.
					// Timeout is not considered an error because it is raised by
					// ReadMessage in absence of messages.
					logger.Printf("Consumer error: %v (%v)\n", err, msg)
					errc <- err
				}
			}
		}()

		logger.Println("Consumer started successfully")

		// Wait for signal.
		logger.Printf("exiting (%v)", <-errc)

		// Send cancellation signal to the goroutines.
		cancel()

		wg.Wait()
		logger.Println("exited")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

type vinDecoderFake struct {
}

func (v vinDecoderFake) Decode(_ context.Context, vinNumber string) (*projection.VINData, error) {
	country := "Spain"
	manufacturer := "Seat"
	brand := "Cupra"
	engineSize := "1200"
	fuelType := "Gasoline"
	model := "Jazz"
	year := "2023"
	assemblyPlant := "Barcelona"
	sn := "411439"
	return &projection.VINData{
		VIN:           vinNumber,
		Country:       &country,
		Manufacturer:  &manufacturer,
		Brand:         &brand,
		EngineSize:    &engineSize,
		FuelType:      &fuelType,
		Model:         &model,
		Year:          &year,
		AssemblyPlant: &assemblyPlant,
		SN:            &sn,
	}, nil
}
