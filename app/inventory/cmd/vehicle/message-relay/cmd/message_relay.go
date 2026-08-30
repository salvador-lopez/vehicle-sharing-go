package cmd

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"vehicle-sharing-go/app/inventory/cmd/vehicle/message-relay/dispatcher/kafka"
	"vehicle-sharing-go/app/inventory/cmd/vehicle/message-relay/relay"
	"vehicle-sharing-go/app/inventory/cmd/vehicle/message-relay/source/binlog"
	binlogredis "vehicle-sharing-go/app/inventory/cmd/vehicle/message-relay/source/binlog/redis"
)

type KafkaTopic struct {
	Name                string
	NumPartitions       int32
	ReplicationFactor   int16
	AggregateTypeRegexp string
}

type KafkaHeaderMappings struct {
	ColumnName string
	HeaderName string
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := canal.NewCanal(getCanalConfig())
		if err != nil {
			return err
		}

		ed, err := getKafkaEventDispatcher()
		if err != nil {
			return err
		}

		stateHandler := getRedisStateHandler()
		handler, posCh, err := relay.NewEventHandler(
			ed,
			viper.GetString("dbAggregateIDColumnName"),
			viper.GetString("dbAggregateTypeColumnName"),
			viper.GetString("dbPayloadColumnName"),
		)
		if err != nil {
			return err
		}

		source := binlog.NewSource(c, handler, posCh, stateHandler, time.Second*5)
		return relay.NewRunner(source).Run()
	},
}

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("includeTransactionTimestamp", true)

	// Configure logrus for relay logging.
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	rootCmd.AddCommand(runCmd)
}

func getKafkaEventDispatcher() (*kafka.EventDispatcher, error) {
	producer, err := getKafkaSyncProducer()
	if err != nil {
		return nil, err
	}

	admin, err := sarama.NewClusterAdmin(viper.GetStringSlice("kafkaBrokers"), sarama.NewConfig())
	if err != nil {
		return nil, err
	}

	var kafkaTopics []KafkaTopic
	if err = viper.UnmarshalKey("kafkaTopics", &kafkaTopics); err != nil {
		return nil, err
	}

	topics := make([]kafka.Topic, 0, len(kafkaTopics))
	for _, topic := range kafkaTopics {
		topics = append(topics, kafka.Topic{
			Name: topic.Name,
			TopicDetail: &sarama.TopicDetail{
				NumPartitions:     topic.NumPartitions,
				ReplicationFactor: topic.ReplicationFactor,
			},
			AggregateType: regexp.MustCompile(topic.AggregateTypeRegexp),
		})
	}

	var kafkaHeaderMappings []kafka.HeaderMapping
	if err = viper.UnmarshalKey("kafkaHeaderMappings", &kafkaHeaderMappings); err != nil {
		return nil, err
	}

	return kafka.NewEventDispatcher(producer, admin, topics, kafkaHeaderMappings)
}

func getKafkaSyncProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	config.Metadata.AllowAutoTopicCreation = false

	return sarama.NewSyncProducer(viper.GetStringSlice("kafkaBrokers"), config)
}

func getRedisStateHandler() *binlogredis.StateHandler {
	return binlogredis.NewStateHandler(
		redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", viper.GetString("redisHost"), viper.GetString("redisPort")),
			DB:   viper.GetInt("redisDB"),
		}),
		viper.GetString("redisKey"),
	)
}

func getCanalConfig() *canal.Config {
	cfg := canal.NewDefaultConfig()

	cfg.Addr = fmt.Sprintf("%s:%s", viper.GetString("dbHost"), viper.GetString("dbPort"))
	cfg.User = viper.GetString("dbUser")
	cfg.Password = viper.GetString("dbPassword")
	cfg.Dump.ExecutionPath = ""
	cfg.IncludeTableRegex = []string{fmt.Sprintf(".*\\.%s", viper.Get("dbOutboxTableRef"))}
	cfg.MaxReconnectAttempts = 10

	return cfg
}
