package bus

import (
	"context"
	"sync"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
)

type kafkaBus[T BusMessage[any]] struct {
	brokers config.ServerAddresses
	topic   string

	tracer trace.Tracer

	mu sync.Mutex
	wg sync.WaitGroup
}

func (k *kafkaBus[T]) Close() {
	client := kafka.Client{
		Addr: kafka.TCP(k.brokers[0].Address()),
	}
	resp, err := client.AlterConfigs(context.Background(), &kafka.AlterConfigsRequest{
		Resources: []kafka.AlterConfigRequestResource{
			{
				ResourceType: kafka.ResourceTypeTopic,
				ResourceName: k.topic,
				Configs: []kafka.AlterConfigRequestConfig{
					{
						Name:  "retention.ms",
						Value: "-1",
					},
				},
			},
		},
		ValidateOnly: false,
	})

	if err != nil {
		log.Logger.Errorf("error updating topic %s: %v", k.topic, err)
		return
	}

	for resource, err := range resp.Errors {

		if err != nil {
			log.Logger.Errorf("error updating topic %s: %s", resource.Name, err.Error())
		}
	}
}
