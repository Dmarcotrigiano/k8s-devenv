package kafka

import (
	"context"
	"fmt"
	"k8s-devenv/kafka/topics"
	"log/slog"
	"os"
	"sync"

	"github.com/hamba/avro/v2"
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

func NewKafka(cfg *KafkaConfig) (*Kafka, error) {
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
		return nil, err
	}

	return &Kafka{
		Client: cl,
		Config: cfg,
		Logger: cfg.Logger,
	}, nil
}

func (k *Kafka) Close() {
	k.Client.Close()
}

func (k *Kafka) Consume(ctx context.Context, handler func(*kgo.Record)) {
	for {
		fetches := k.Client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return
		}
		fetches.EachError(func(t string, p int32, err error) {
			k.Logger.Error("kafka fetch error", "topic", t, "partition", p, "error", err)
			os.Exit(1)
		})
		fetches.EachRecord(handler)
		if err := k.Client.CommitUncommittedOffsets(ctx); err != nil {
			k.Logger.Error("commit records failed", "error", err)
			continue
		}
	}
}

func (k *Kafka) ProduceAvroMessage(ctx context.Context, topic topics.KafkaTopic, data interface{}) error {
	schema := topic.Struct().Schema()
	k.Logger.Debug("got Avro schema", "topic", topic, "schema", schema.Fingerprint())
	bytes, err := avro.Marshal(schema, data)
	if err != nil {
		k.Logger.Debug("Error marshalling data", "topic", topic.String(), "data", data, "error", err)
		return fmt.Errorf("error marshalling data: %w", err)
	}
	record := &kgo.Record{
		Topic: topic.String(),
		Value: bytes,
	}
	k.Logger.Debug("Producing message", "topic", topic, "data", data, "length", len(bytes), "record", record, "bytes", bytes)
	if err := k.Client.ProduceSync(ctx, record).FirstErr(); err != nil {
		k.Logger.Debug("Error producing message", "error", err)
		return fmt.Errorf("error producing message: %w", err)
	}
	return nil
}

func (k *Kafka) ProduceAvroMessages(ctx context.Context, topic topics.KafkaTopic, data []interface{}) error {
	var wg sync.WaitGroup
	for _, d := range data {
		wg.Add(1)
		go func(d interface{}) {
			defer wg.Done()
			schema := topic.Struct().Schema()
			k.Logger.Debug("got Avro schema", "topic", topic, "schema", schema.Fingerprint())
			bytes, err := avro.Marshal(schema, d)
			if err != nil {
				k.Logger.Debug("Error marshalling data", "topic", topic.String(), "data", d, "error", err)
				return
			}
			record := &kgo.Record{
				Topic: topic.String(),
				Value: bytes,
			}
			k.Logger.Debug("Producing message", "topic", topic, "data", d, "length", len(bytes), "record", record, "bytes", bytes)
			k.Client.Produce(ctx, record, func(_ *kgo.Record, err error) {
				if err != nil {
					k.Logger.Error("record had a produce error", "error", err)
				}
			})
		}(d)
	}
	wg.Wait()
	return nil
}
