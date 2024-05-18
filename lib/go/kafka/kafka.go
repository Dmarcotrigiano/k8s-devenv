package kafka

import "github.com/twmb/franz-go/pkg/kgo"

type Kafka struct {
	Client *kgo.Client
}

func New(opts ...kgo.Opt) *Kafka {
	cl, err := kgo.NewClient(opts...)
	if err != nil {
		panic(err)
	}

	return &Kafka{
		Client: cl,
	}
}

func (k *Kafka) Close() {
	k.Client.Close()
}

func (k *Kafka) Produce(msg) error {

}
