package config

import "k8s-devenv/kafka"

type Config struct {
	Kafka *kafka.Kafka
	Scylla
}
