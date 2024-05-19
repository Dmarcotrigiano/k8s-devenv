package kafka

import (
	"k8s-devenv/kafka/topics"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Kafka struct {
	Client *kgo.Client
	Config *KafkaConfig
	Logger *slog.Logger
}

type KafkaConfig struct {
	ConsumeTopics []topics.KafkaTopic
	Brokers       []string
	ClientID      string
	ConsumerGroup string
	Logger        *slog.Logger
}

func New(cfg *KafkaConfig) *Kafka {
	consumeTopics := make([]string, len(cfg.ConsumeTopics))
	for i, topic := range cfg.ConsumeTopics {
		consumeTopics[i] = topic.String()
	}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumeTopics(consumeTopics...),
		kgo.ClientID(cfg.ClientID),
		kgo.ConsumerGroup(cfg.ConsumerGroup),
		// If you wish to add more configuration options, you can do so here.
	)
	if err != nil {
		panic(err)
	}

	return &Kafka{
		Client: cl,
		Config: cfg,
		Logger: cfg.Logger,
	}
}

func (k *Kafka) Close() {
	k.Client.Close()
}
