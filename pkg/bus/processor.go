package bus

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ShatteredRealms/go-common-service/pkg/log"
)

type BusProcessor interface {
	StartProcessing(ctx context.Context)
	StopProcessing()
}

var (
	ErrFetchMessage     = errors.New("unable to fetch message")
	ErrProcessingFailed = errors.New("processing failed")
)

type DefaultBusProcessor[T BusModelMessage[any]] struct {
	Reader             MessageBusReader[T]
	Repo               BusMessageRepository[T]
	mu                 sync.Mutex
	concurrentFetchErr int
	concurrentErrCount int
	isProcessing       bool
}

func (bp *DefaultBusProcessor[T]) StartProcessing(ctx context.Context) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	if bp.isProcessing {
		return
	}

	bp.isProcessing = true

	go func() {
		defer func() {
			bp.mu.Lock()
			bp.isProcessing = false
			bp.mu.Unlock()
		}()

		for bp.isProcessing {
			err := bp.process(ctx)
			if errors.Is(err, ErrProcessingFailed) {
				bp.concurrentErrCount++
			} else if errors.Is(err, ErrFetchMessage) {
				bp.concurrentFetchErr++
			} else {
				bp.concurrentErrCount = 0
				bp.concurrentFetchErr = 0
			}

			if bp.concurrentErrCount >= 5 {
				log.Logger.WithContext(ctx).Error("too many errors processing dimension messages")
				return
			}

			if bp.concurrentFetchErr >= 10 {
				log.Logger.WithContext(ctx).Error("too many errors fetching dimension messages")
			}
		}
	}()
}

func (bp *DefaultBusProcessor[T]) process(ctx context.Context) error {
	msg, err := bp.Reader.FetchMessage(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("unable to fetch %T message: %v", msg, err)
		bp.Reader.ProcessFailed()
		return fmt.Errorf("%w: %w", ErrFetchMessage, err)
	}

	if (*msg).WasDeleted() {
		err = bp.Repo.Delete(ctx, (*msg).GetId())
		if err != nil {
			log.Logger.WithContext(ctx).Errorf(
				"unable to delete %T %s: %v", msg, (*msg).GetId(), err)
			bp.Reader.ProcessFailed()
			return ErrProcessingFailed
		}

		bp.Reader.ProcessSucceeded(ctx)
		return nil
	}

	err = bp.Repo.Save(ctx, *msg)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf(
			"unable to save %T %s: %v", msg, (*msg).GetId(), err)
		bp.Reader.ProcessFailed()
		return ErrProcessingFailed
	}

	bp.Reader.ProcessSucceeded(ctx)
	return nil
}

func (bp *DefaultBusProcessor[T]) StopProcessing() {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.isProcessing = false
}
