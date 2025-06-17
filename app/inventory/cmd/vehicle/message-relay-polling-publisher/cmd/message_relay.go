package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	gormpkg "vehicle-sharing-go/pkg/database/gorm"
	"vehicle-sharing-go/pkg/messaging"
	kafkapkg "vehicle-sharing-go/pkg/messaging/kafka"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers []string
	GroupId string
	Topic   string
}

type DbConfig struct {
	Conn DbConn
	Name string
}

type DbConn struct {
	Host     string
	Port     int
	User     string
	Password string
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Message Relay Polling Publisher",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		errc := make(chan error)

		// Graceful shutdown
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			errc <- fmt.Errorf("%s", <-c)
		}()

		logger := log.New(os.Stderr, "[inventory-vehicles-message-relay-polling-publisher] ", log.Ltime)

		var dbCfg DbConfig
		err := viper.UnmarshalKey("db", &dbCfg)
		if err != nil {
			logger.Fatal(err)
		}
		connCfg := dbCfg.Conn
		dbConn, err := gormpkg.NewConnectionFromConfig(&gormpkg.Config{
			Host:         connCfg.Host,
			Port:         connCfg.Port,
			UserName:     connCfg.User,
			Password:     connCfg.Password,
			DatabaseName: dbCfg.Name,
			Logger:       logger,
			LogQueries:   false,
		})
		if err != nil {
			logger.Fatalf("failed to create db connection: %v", err)
		}

		outboxRepo := gormpkg.NewOutboxRepository(dbConn)

		// Load config
		var kafkaCfg KafkaConfig
		err = viper.UnmarshalKey("kafka", &kafkaCfg)
		if err != nil {
			logger.Fatal(err)
		}

		producer, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers":     strings.Join(kafkaCfg.Brokers, ","),
			"broker.address.family": "v4",
			"group.id":              kafkaCfg.GroupId,
			"auto.offset.reset":     "earliest",
		})
		if err != nil {
			logger.Fatalf("failed to create kafka producer: %v", err)
		}
		defer producer.Close()

		eventPublisher := kafkapkg.NewEventPublisher(producer, kafkaCfg.Topic)

		// Poll loop
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		p := messaging.NewPollingPublisher(
			outboxRepo,
			eventPublisher,
			2*time.Second,
			10,
			logger,
		)

		go func() {
			if err := p.Run(ctx); err != nil {
				errc <- err
			}
		}()

		logger.Println("Massage Relay Polling Publisher started successfully")

		// Wait for error signal.
		logger.Printf("exiting (%v)", <-errc)

		return nil
	},
}

func init() {
	viper.AutomaticEnv()

	rootCmd.AddCommand(runCmd)
}
