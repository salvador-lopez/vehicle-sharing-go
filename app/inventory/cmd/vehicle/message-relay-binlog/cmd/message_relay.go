package cmd

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-redis/redis/v8"
	"github.com/lorenzoranucci/tor/adapters/kafka"
	redisadapter "github.com/lorenzoranucci/tor/adapters/redis"
	"github.com/lorenzoranucci/tor/router/pkg/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers        []string
	Topics         []KafkaTopic
	HeaderMappings []kafka.HeaderMapping
}
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

type DbConfig struct {
	Conn                    DbConn
	Name                    string
	OutboxTableRef          string
	PayloadColumnName       string
	AggregateIdColumnName   string
	AggregateTypeColumnName string
}

type DbConn struct {
	Host     string
	User     string
	Password string
	Port     string
}

type RedisConfig struct {
	Host string
	Port string
	Key  string
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	RunE: func(cmd *cobra.Command, args []string) error {
		var dbCfg DbConfig
		err := viper.UnmarshalKey("db", &dbCfg)
		if err != nil {
			return err
		}

		c, err := canal.NewCanal(getCanalConfig(&dbCfg))
		if err != nil {
			return err
		}

		ed, err := getKafkaEventDispatcher()
		if err != nil {
			return err
		}

		stateHandler, err := getRedisStateHandler()
		if err != nil {
			return err
		}

		handler, err := run.NewEventHandler(
			ed,
			dbCfg.AggregateIdColumnName,
			dbCfg.AggregateTypeColumnName,
			dbCfg.PayloadColumnName,
		)
		if err != nil {
			return err
		}

		runner := run.NewRunner(c, handler, stateHandler, time.Second*5)

		return runner.Run()
	},
}

func init() {
	viper.AutomaticEnv()

	rootCmd.AddCommand(runCmd)
}

func getKafkaEventDispatcher() (*kafka.EventDispatcher, error) {
	var kafkaCfg KafkaConfig
	err := viper.UnmarshalKey("kafka", &kafkaCfg)
	if err != nil {
		return nil, err
	}

	producer, err := getKafkaSyncProducer(kafkaCfg.Brokers)
	if err != nil {
		return nil, err
	}

	admin, err := sarama.NewClusterAdmin(kafkaCfg.Brokers, sarama.NewConfig())
	if err != nil {
		return nil, err
	}

	topics := make([]kafka.Topic, 0, len(kafkaCfg.Topics))
	for _, topic := range kafkaCfg.Topics {
		topics = append(topics, kafka.Topic{
			Name: topic.Name,
			TopicDetail: &sarama.TopicDetail{
				NumPartitions:     topic.NumPartitions,
				ReplicationFactor: topic.ReplicationFactor,
			},
			AggregateType: regexp.MustCompile(topic.AggregateTypeRegexp),
		})
	}

	return kafka.NewEventDispatcher(producer, admin, topics, kafkaCfg.HeaderMappings)
}

func getKafkaSyncProducer(kafkaBrokersCfg []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	config.Metadata.AllowAutoTopicCreation = false

	producer, err := sarama.NewSyncProducer(kafkaBrokersCfg, config)
	if err != nil {
		return nil, err
	}

	return producer, err
}

func getRedisStateHandler() (*redisadapter.StateHandler, error) {
	var redisCfg RedisConfig
	err := viper.UnmarshalKey("redis", &redisCfg)
	if err != nil {
		return nil, err
	}

	return redisadapter.NewStateHandler(
		redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", redisCfg.Host, redisCfg.Port)}),
		redisCfg.Key,
	), nil
}

func getCanalConfig(dbCfg *DbConfig) *canal.Config {
	dbConn := dbCfg.Conn

	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%s", dbConn.Host, dbConn.Port)
	cfg.User = dbConn.User
	cfg.Password = dbConn.Password
	cfg.Dump.ExecutionPath = ""
	cfg.IncludeTableRegex = []string{fmt.Sprintf(".*\\.%s", dbCfg.OutboxTableRef)}
	cfg.MaxReconnectAttempts = 10

	return cfg
}
