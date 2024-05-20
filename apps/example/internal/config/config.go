package config

import (
	"k8s-devenv/kafka"
	"k8s-devenv/kafka/topics"
	"k8s-devenv/scylla"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Kafka  *kafka.Kafka
	Scylla *scylla.Cluster
	Logger *slog.Logger
}

func Load() *Config {
	isDebug := os.Getenv("DEBUG")
	logLevel := slog.LevelInfo
	if isDebug == "true" {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	logger.Info("Starting example-1 service")
	exampleOneTopic := topics.TopicOne
	kafkaConfig := &kafka.KafkaConfig{
		ConsumeTopics: []topics.KafkaTopic{exampleOneTopic},
	}
	kc, err := kafka.NewKafka(kafkaConfig)
	if err != nil {
		logger.Error("Error creating Kafka client", "error", err)
		os.Exit(1)
	}
	logger.Debug("Kafka client created")
	logger.Debug("Initiating connection to Scylla")
	scyllaHosts := os.Getenv("SCYLLA_HOST")
	ks := os.Getenv("SCYLLA_KEYSPACE")
	if scyllaHosts == "" || ks == "" {
		logger.Error("SCYLLA_HOST or SCYLLA_KEYSPACE not set", "SCYLLA_HOST", scyllaHosts, "SCYLLA_KEYSPACE", ks)
		os.Exit(1)
	}
	logger.Debug("Creating Scylla session", scyllaHosts, ks)
	sc := scylla.NewScyllaCluster(strings.Split(scyllaHosts, " "), ks)
	err = sc.CreateSession()
	if err != nil {
		logger.Error("error creating Scylla session", "error", err)
		os.Exit(1)
	}
	logger.Debug("Scylla session created")
	logger.Debug("Creating keyspace and tables")
	return &Config{
		Kafka:  kc,
		Scylla: sc,
		Logger: logger,
	}
}
