package topics

import "github.com/hamba/avro/v2"

type KafkaTopic int

type GeneratedStruct interface {
	Schema() avro.Schema
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
}

const (
	TopicOne KafkaTopic = iota
	TopicTwo
)

func (kt KafkaTopic) String() string {
	return [...]string{"example-topic-one", "example-topic-two"}[kt]
}

func (kt KafkaTopic) Struct() interface{} {
	return [...]GeneratedStruct{&ExampleTopicOne{}, &ExampleTopic2{}}[kt]
}
