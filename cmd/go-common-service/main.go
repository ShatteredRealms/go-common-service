package main

import (
	"context"
	"errors"
	"io"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/character/characterbus"
	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/go-faker/faker/v4"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	resetBus     = true
	sendMessages = false
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
				inCtx := context.Background()
				msg, err := b.FetchMessage(inCtx)
				if err != nil {
					if !errors.Is(err, io.EOF) {
						log.Logger.Errorf("Error for group %s fetching message: %v", b.GetGroup(), err)
					}
					continue
				}

				if failCount < maxFailCount {
					failCount++
					log.Logger.Infof("Failing for group %s to process and message: %v", b.GetGroup(), msg)
					err := b.ProcessFailed()
					if err != nil {
						log.Logger.Errorf("Failed for group %s to mark %v as failed: %v", b.GetGroup(), msg, err)
					}
					continue
				}

				failCount = 0
				log.Logger.Infof("Succeeding for group %s to process message: %v", b.GetGroup(), msg)
				err = b.ProcessSucceeded(inCtx)
				if err != nil {
					log.Logger.Errorf("Failed for group %s to mark %v as succeeded: %v", b.GetGroup(), msg, err)
				}
			}
		}()
	}

	if resetBus {
		log.Logger.Info("Resetting bus")
		bus := readBusses[0]
		err := bus.Reset(ctx)
		if err != nil {
			log.Logger.Errorf("Error resetting bus: %v", err)
		}
		log.Logger.Info("Resetting bus complete")
	}

	writeBus := bus.NewKafkaMessageBusWriter([]config.ServerAddress{{Host: "localhost", Port: "29092"}}, msg)
	ticker := time.NewTicker(5 * time.Second)
	tracer := tp.Tracer("main")
	if !sendMessages {
		ticker.Stop()
	}
	for {
		select {
		case <-ticker.C:
			ctx, span := tracer.Start(ctx, "publish-message")
			newMsg := characterbus.Message{
				Id:          faker.UUIDHyphenated(),
				OwnerId:     faker.UUIDHyphenated(),
				DimensionId: faker.UUIDHyphenated(),
				MapId:       faker.UUIDHyphenated(),
				Deleted:     false,
			}

			log.Logger.Infof("Publishing message (%s)", newMsg.GetId())
			writeBus.Publish(ctx, newMsg)
			span.End()

		case <-ctx.Done():
			log.Logger.Info("Shut down requested by user")

			wg := sync.WaitGroup{}
			for _, b := range readBusses {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := b.Close()
					if err != nil {
						log.Logger.Errorf("Error shutting down bus: %v", err)
					}
					log.Logger.Infof("Bus %s shut down complete", b.GetMessageType())
				}()
			}

			wg.Wait()
			log.Logger.Info("Shut down complete")
			return
		}
	}
}
