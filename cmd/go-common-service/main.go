package main

import (
	"context"
	"errors"
	"io"
	"os"
	"os/signal"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/character/characterbus"
	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	log.Logger.Level = logrus.InfoLevel
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	log.Logger.Info("Starting program")
	tp := trace.NewTracerProvider()
	defer tp.Shutdown(ctx)

	msg := characterbus.Message{}
	readBusses := make([]bus.MessageBusReader[characterbus.Message], 0)
	cg1, cg2 := "service1", "service2"
	readBusses = append(readBusses, bus.NewKafkaMessageBusReader([]config.ServerAddress{{Host: "localhost", Port: "29092"}}, cg1, msg))
	readBusses = append(readBusses, bus.NewKafkaMessageBusReader([]config.ServerAddress{{Host: "localhost", Port: "29092"}}, cg1, msg))
	readBusses = append(readBusses, bus.NewKafkaMessageBusReader([]config.ServerAddress{{Host: "localhost", Port: "29092"}}, cg2, msg))
	readBusses = append(readBusses, bus.NewKafkaMessageBusReader([]config.ServerAddress{{Host: "localhost", Port: "29092"}}, cg2, msg))

	for _, b := range readBusses {
		go func() {
			failCount := 0
			maxFailCount := 0
			for ctx.Err() == nil {
				msg, err := b.FetchMessage(ctx)
				if err != nil {
					if !errors.Is(err, io.EOF) {
						log.Logger.Errorf("Error fetching message: %v", err)
					}
					continue
				}

				if failCount < maxFailCount {
					failCount++
					log.Logger.Infof("Failing to process message: %v", msg)
					err := b.ProcessFailed()
					if err != nil {
						log.Logger.Errorf("Failed to mark %v as failed: %v", msg, err)
					}
					continue
				}

				failCount = 0
				log.Logger.Infof("Succeeding to process message: %v", msg)
				err = b.ProcessSucceeded(ctx)
				if err != nil {
					log.Logger.Errorf("Failed to mark %v as succeeded: %v", msg, err)
				}
			}
		}()
	}

	// writeBus := bus.NewKafkaMessageBusWriter([]config.ServerAddress{{Host: "localhost", Port: "29092"}}, msg)
	// ticker := time.NewTicker(5 * time.Second)
	// tracer := tp.Tracer("main")
	for {
		select {
		// case <-ticker.C:
		// 	ctx, span := tracer.Start(ctx, "publish-message")
		// 	newMsg := bus.CharacterCreatedMessage{ID: faker.UUIDHyphenated()}
		//
		// 	log.Logger.Infof("Publishing message (%s)", newMsg.GetId())
		// 	writeBus.Publish(ctx, newMsg)
		// 	span.End()

		case <-ctx.Done():
			log.Logger.Info("Shut down requested by user")
			for _, b := range readBusses {
				err := b.Close()
				if err != nil {
					log.Logger.Errorf("Error shutting down bus: %v", err)
				}
			}
			log.Logger.Info("Shut down complete")
			return
		}
	}
}
