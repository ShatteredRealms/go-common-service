package mapbus

import (
	"context"
	"sync"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
)

type Service interface {
	GetMaps(ctx context.Context) (*Maps, error)
	GetMapById(ctx context.Context, mId string) (*Map, error)
	StartProcessingBus(ctx context.Context)
	StopProcessingBus()
}

type service struct {
	repo               Repository
	mBus               bus.MessageBusReader[Message]
	isProcessing       bool
	concurrentErrCount int

	mu sync.Mutex
}

func NewService(
	repo Repository,
	mBus bus.MessageBusReader[Message],
) Service {
	return &service{
		repo: repo,
		mBus: mBus,
	}
}

// StartProcessingBus implements MapService.
func (d *service) StartProcessingBus(ctx context.Context) {
	d.mu.Lock()
	if d.isProcessing {
		d.mu.Unlock()
		return
	}

	d.mu.Lock()
	d.isProcessing = true
	d.mu.Unlock()

	go func() {
		for d.isProcessing && d.concurrentErrCount < 5 {
			msg, err := d.mBus.FetchMessage(ctx)
			if err != nil {
				log.Logger.WithContext(ctx).Errorf("unable to fetch map message: %v", err)
				continue
			}

			if msg.Deleted {
				_, err = d.repo.DeleteMap(ctx, msg.Id)
				if err != nil {
					log.Logger.WithContext(ctx).Errorf(
						"unable to delete map %s: %v", msg.Id, err)
					d.mBus.ProcessFailed()
					d.concurrentErrCount++
				} else {
					d.mBus.ProcessSucceeded(ctx)
					d.concurrentErrCount = 0
				}
			} else {
				d.repo.CreateMap(ctx, msg.Id)
				if err != nil {
					log.Logger.WithContext(ctx).Errorf(
						"unable to save map %s: %v", msg.Id, err)
					d.mBus.ProcessFailed()
					d.concurrentErrCount++
				} else {
					d.mBus.ProcessSucceeded(ctx)
					d.concurrentErrCount = 0
				}
			}
		}
	}()
}

// StopProcessingBus implements MapService.
func (d *service) StopProcessingBus() {
	d.isProcessing = false
}

// GetMapById implements MapService.
func (d *service) GetMapById(ctx context.Context, mId string) (*Map, error) {
	return d.repo.GetMapById(ctx, mId)
}

// GetMaps implements MapService.
func (d *service) GetMaps(ctx context.Context) (*Maps, error) {
	return d.repo.GetMaps(ctx)
}
