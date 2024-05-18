package kafka

import (
	"testing"
)

func TestKafka(t *testing.T) {
	result := Kafka("works")
	if result != "Kafka works" {
		t.Error("Expected Kafka to append 'works'")
	}
}
