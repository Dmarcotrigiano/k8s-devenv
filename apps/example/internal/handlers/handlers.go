package handlers

import (
	"k8s-devenv/example/internal/config"
	"k8s-devenv/kafka/topics"
	"k8s-devenv/scylla/example"
	"k8s-devenv/scylla/models"
	"os"

	"github.com/google/uuid"
	"github.com/hamba/avro/v2"
	"github.com/twmb/franz-go/pkg/kgo"
)

func GetMessageHandler(cfg *config.Config) func(msg *kgo.Record) {
	cfg.Logger.Debug("Creating message handler")
	topicOne := topics.TopicOne.String()
	return func(msg *kgo.Record) {
		switch msg.Topic {
		case topicOne:
			to := topics.TopicOne.Struct().(*topics.ExampleTopicOne)
			err := avro.Unmarshal(to.Schema(), msg.Value, &to)
			if err != nil {
				cfg.Logger.Error("error unmarshalling message", "error", err)
				os.Exit(1)
			}
			cfg.Logger.Debug("Unmarshaled data", "data", to, "raw", msg.Value)
			handleTopicOneMessage(cfg, to)
		default:
			cfg.Logger.Error("unknown topic", "topic", msg.Topic)
			os.Exit(1)
		}
	}
}

func handleTopicOneMessage(cfg *config.Config, to *topics.ExampleTopicOne) {
	cfg.Logger.Debug("Handling topic one message", "message", to)
	example.InsertExample(cfg.Scylla, &models.Example{
		ID: uuid.New(),
		A:  to.A,
		B:  to.B,
	})
}
