package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/go-faker/faker/v4"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	log.Logger.Info("Starting program")
	tp := trace.NewTracerProvider()
	defer tp.Shutdown(ctx)

	tracer := tp.Tracer("main")
	msg := bus.CharacterCreatedMessage{}
	busses := make([]bus.MessageBus[bus.CharacterCreatedMessage], 0)
	cg1, cg2 := "service1", "service2"
	busses = append(busses, newMessageBus(ctx, cg1, msg))
	busses = append(busses, newMessageBus(ctx, cg1, msg))
	busses = append(busses, newMessageBus(ctx, cg2, msg))
	busses = append(busses, newMessageBus(ctx, cg2, msg))

	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			ctx, span := tracer.Start(ctx, "publish-message")
			newMsg := bus.CharacterCreatedMessage{ID: faker.UUIDHyphenated(), Name: faker.Username()}

			log.Logger.Infof("Publishing message (%s)", newMsg.GetId())
			busses[0].Publish(ctx, newMsg)
			span.End()

		case <-ctx.Done():
			log.Logger.Info("Shut down requested by user")
			for _, b := range busses {
				err := b.Close(ctx)
				if err != nil {
					log.Logger.Errorf("Error shutting down bus: %v", err)
				}
			}
			log.Logger.Info("Shut down complete")
			return
		}
	}
}

func newMessageBus[T any](ctx context.Context, svc string, msg bus.BusMessage[T]) bus.MessageBus[T] {
	b := bus.NewKafkaMessageBus([]config.ServerAddress{{Host: "localhost", Port: "29092"}}, svc, msg)
	c := make(chan T)
	go func() {
		log.Logger.Infof("Listening for messages for group.id=%s", svc)
		b.ReceiveMessages(ctx, c)
	}()
	go func() {
		for {
			select {
			case msg := <-c:
				log.Logger.Infof("Received message for group.id=%s: %v", svc, msg)
			}
		}
	}()
	return b
}
