package bus

import (
	"errors"
	"sync"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrSerializeMessage = errors.New("unable to serialize message")
	ErrSendingMessage   = errors.New("unable to send message on bus")
	ErrEncodingMessage  = errors.New("unable to encode message")
)

type kafkaBus[T BusMessage[any]] struct {
	brokers config.ServerAddresses
	topic   string

	tracer trace.Tracer

	mu sync.Mutex
	wg sync.WaitGroup
}
