package dimensionbus

import (
	"context"
	"sync"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
)

type Service interface {
	GetDimensions(ctx context.Context) (*Dimensions, error)
	GetDimensionById(ctx context.Context, dimensionId string) (*Dimension, error)
	StartProcessingBus(ctx context.Context)
	StopProcessingBus()
}

type service struct {
	repo               Repository
	dimensionBus       bus.MessageBusReader[Message]
	isProcessing       bool
	concurrentErrCount int

	mu sync.Mutex
}

func NewService(
	repo Repository,
	dimensionBus bus.MessageBusReader[Message],
) Service {
	return &service{
		repo:         repo,
		dimensionBus: dimensionBus,
	}
}

// StartProcessingBus implements DimensionService.
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
			msg, err := d.dimensionBus.FetchMessage(ctx)
			if err != nil {
				log.Logger.WithContext(ctx).Errorf("unable to fetch dimension message: %v", err)
				continue
			}

			if msg.Deleted {
				_, err = d.repo.DeleteDimension(ctx, msg.Id)
				if err != nil {
					log.Logger.WithContext(ctx).Errorf(
						"unable to delete dimension %s: %v", msg.Id, err)
					d.dimensionBus.ProcessFailed()
					d.concurrentErrCount++
				} else {
					d.dimensionBus.ProcessSucceeded(ctx)
					d.concurrentErrCount = 0
				}
			} else {
				d.repo.CreateDimension(ctx, msg.Id)
				if err != nil {
					log.Logger.WithContext(ctx).Errorf(
						"unable to save dimension %s: %v", msg.Id, err)
					d.dimensionBus.ProcessFailed()
					d.concurrentErrCount++
				} else {
					d.dimensionBus.ProcessSucceeded(ctx)
					d.concurrentErrCount = 0
				}
			}
		}
	}()
}

// StopProcessingBus implements DimensionService.
func (d *service) StopProcessingBus() {
	d.isProcessing = false
}

// GetDimensionById implements DimensionService.
func (d *service) GetDimensionById(ctx context.Context, dimensionId string) (*Dimension, error) {
	return d.repo.GetDimensionById(ctx, dimensionId)
}

// GetDimensions implements DimensionService.
func (d *service) GetDimensions(ctx context.Context) (*Dimensions, error) {
	return d.repo.GetDimensions(ctx)
}
