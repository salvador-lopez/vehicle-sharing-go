package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
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
			UserName:     "root",
			Password:     "root",
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

		// A signal handler or similar could be used to set this to false to break the loop.
		run := true

		for run {
			msg, err := c.ReadMessage(time.Second)
			if err == nil {
				logger.Printf("Message on %s with AggregateID %s: %s\n", msg.TopicPartition, string(msg.Key), string(msg.Value))

				aggregateID, err := uuid.Parse(string(msg.Key))
				if err != nil {
					logger.Fatalf("Error Unmarshalling message aggregateID: %v (%v)\n", err, msg)
				}

				var payload *event.CarCreatedPayload
				err = json.Unmarshal(msg.Value, &payload)
				if err != nil {
					logger.Fatalf("Error Unmarshalling message: %v (%v)\n", err, msg)
				}

				err = carProjector.ProjectCarCreated(ctx, aggregateID, payload)
				if err != nil {
					logger.Fatalf("Error Projecting Event %s: %v (%v)\n", aggregateID, err, payload)
				}
			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				logger.Fatalf("Consumer error: %v (%v)\n", err, msg)

			}
		}

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

func toTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}

func Decode(input any, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:  "json",
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			toTimeHookFunc()),
		Result: result,
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}
