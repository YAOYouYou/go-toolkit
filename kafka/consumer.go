package kafkautil

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/YAOYouYou/go-toolkit/logging"
)

var logger = logging.GetDefaultLogger()

type ConsumerConfig struct {
	Topic            string
	GroupId          string
	BootstrapServers string

	HearBeatIntervalMs     int
	SessionTimeoutMs       int
	MaxPollIntervalMs      int
	FetchMaxBytes          int
	MaxPartitionFetchBytes int
	SecurityProtocol       string
	SaslMechanism          string
	SaslUsername           string
	SaslPassword           string
}

func NewConsumer(cfg *ConsumerConfig) *kafka.Consumer {
	var kafkaconf = &kafka.ConfigMap{
		"api.version.request":       "true",
		"auto.offset.reset":         "latest",
		"heartbeat.interval.ms":     3000,
		"session.timeout.ms":        30000,
		"max.poll.interval.ms":      120000,
		"fetch.max.bytes":           1024000,
		"max.partition.fetch.bytes": 256000}
	kafkaconf.SetKey("security.protocol", "plaintext")
	kafkaconf.SetKey("bootstrap.servers", cfg.BootstrapServers)
	kafkaconf.SetKey("group.id", cfg.GroupId)
	if cfg.HearBeatIntervalMs != 0 {
		kafkaconf.SetKey("heartbeat.interval.ms", cfg.HearBeatIntervalMs)
	}
	if cfg.SessionTimeoutMs != 0 {
		kafkaconf.SetKey("session.timeout.ms", cfg.SessionTimeoutMs)
	}
	if cfg.MaxPollIntervalMs != 0 {
		kafkaconf.SetKey("max.poll.interval.ms", cfg.MaxPollIntervalMs)
	}
	if cfg.FetchMaxBytes != 0 {
		kafkaconf.SetKey("fetch.max.bytes", cfg.FetchMaxBytes)
	}
	if cfg.MaxPartitionFetchBytes != 0 {
		kafkaconf.SetKey("max.partition.fetch.bytes", cfg.MaxPartitionFetchBytes)
	}

	// switch cfg.SecurityProtocol {
	// case "plaintext":
	// 	kafkaconf.SetKey("security.protocol", "plaintext")
	// case "sasl_ssl":
	// 	kafkaconf.SetKey("security.protocol", "sasl_ssl")
	// 	kafkaconf.SetKey("ssl.ca.location", "./conf/ca-cert.pem")
	// 	kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
	// 	kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
	// 	kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
	// case "sasl_plaintext":
	// 	kafkaconf.SetKey("security.protocol", "sasl_plaintext")
	// 	kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
	// 	kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
	// 	kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)

	// default:
	// 	panic(kafka.NewError(kafka.ErrUnknownProtocol, "unknown protocol", true))
	// }

	consumer, err := kafka.NewConsumer(kafkaconf)
	if err != nil {
		panic(err)
	}
	logger.Debugf("init kafka consumer success\n")
	return consumer
}
