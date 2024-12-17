package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/character/characterbus"
	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	resetBus     = true
	sendMessages = true
)

type IgnoreRepo struct{}

func (i *IgnoreRepo) Save(ctx context.Context, data characterbus.Message) error {
	return nil
}
func (i *IgnoreRepo) Delete(ctx context.Context, id *uuid.UUID) error {
	return nil
}

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

	processors := make([]bus.BusProcessor[characterbus.Message], 0)

	for _, b := range readBusses {
		busProcesor := bus.DefaultBusProcessor[characterbus.Message]{
			Reader: b,
			Repo:   &IgnoreRepo{},
		}
		busProcesor.StartProcessing(ctx)
		processors = append(processors, &busProcesor)
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
			id, err := uuid.NewV7()
			if err != nil {
				log.Logger.Errorf("Error generating UUID: %v", err)
				break
			}
			oId, err := uuid.NewV7()
			if err != nil {
				log.Logger.Errorf("Error generating UUID: %v", err)
				break
			}
			dId, err := uuid.NewV7()
			if err != nil {
				log.Logger.Errorf("Error generating UUID: %v", err)
				break
			}
			mId, err := uuid.NewV7()
			if err != nil {
				log.Logger.Errorf("Error generating UUID: %v", err)
				break
			}
			newMsg := characterbus.Message{
				Id:          id,
				OwnerId:     oId,
				DimensionId: dId,
				MapId:       mId,
				Deleted:     false,
			}

			log.Logger.Infof("Publishing message (%s)", newMsg.GetId())
			writeBus.Publish(ctx, newMsg)
			span.End()

		case <-ctx.Done():
			log.Logger.Info("Shut down requested by user")

			wg := sync.WaitGroup{}
			for _, p := range processors {
				wg.Add(1)
				go func() {
					defer wg.Done()
					p.StopProcessing()
				}()
			}
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
