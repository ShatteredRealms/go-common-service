package bus

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/google/uuid"
)

type BusListener[T BusMessage[any]] func(ctx context.Context, msg *T)
type BusListenerHandler int

type BusProcessor[T BusMessage[any]] interface {
	StartProcessing(ctx context.Context)
	StopProcessing()
	IsProcessing() bool
	GetResetter() Resettable
	RegisterListener(listener BusListener[T]) BusListenerHandler
	RemoveListener(listenerHandle BusListenerHandler)
}

var (
	ErrFetchMessage     = errors.New("unable to fetch message")
	ErrProcessingFailed = errors.New("processing failed")
)

type DefaultBusProcessor[T BusModelMessage[any]] struct {
	Reader             MessageBusReader[T]
	Repo               BusMessageRepository[T]
	Listeners          map[BusListenerHandler]BusListener[T]
	nextListenerHandle BusListenerHandler
	mu                 sync.Mutex
	concurrentFetchErr int
	concurrentErrCount int
	isProcessing       bool
}

// GetResetter implements BusProcessor.
func (bp *DefaultBusProcessor[T]) GetResetter() Resettable {
	return bp.Reader
}

func (bp *DefaultBusProcessor[T]) IsProcessing() bool {
	return bp.isProcessing
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

		log.Logger.WithContext(ctx).Infof("Starting bus processor for %s", bp.Reader.GetMessageType())
		for bp.isProcessing {
			msg, err := bp.process(ctx)
			if errors.Is(err, ErrProcessingFailed) {
				bp.concurrentErrCount++
			} else if errors.Is(err, ErrFetchMessage) {
				if errors.Is(err, context.Canceled) {
					log.Logger.WithContext(ctx).Info("stopping bus processor due to context cancellation")
					return
				}
				bp.concurrentFetchErr++
			} else {
				bp.concurrentErrCount = 0
				bp.concurrentFetchErr = 0

				if msg != nil {
					for _, listener := range bp.Listeners {
						listener(ctx, msg)
					}
				}
			}

			if bp.concurrentErrCount >= 5 {
				log.Logger.WithContext(ctx).Error("too many errors processing dimension messages")
				return
			}

			if bp.concurrentFetchErr >= 10 {
				log.Logger.WithContext(ctx).Error("too many errors fetching dimension messages")
				return
			}
		}
	}()
}

func (bp *DefaultBusProcessor[T]) process(ctx context.Context) (*T, error) {
	msg, err := bp.Reader.FetchMessage(ctx)
	if err != nil {
		if errors.Is(err, ErrDecodingMessage) {
			log.Logger.Warnf("skipping %T message due to invalid format", msg)
			bp.Reader.ProcessSkipped(ctx)
			return nil, nil
		}

		log.Logger.WithContext(ctx).Errorf("unable to fetch %T message: %v", msg, err)
		bp.Reader.ProcessFailed()
		return nil, fmt.Errorf("%w: %w", ErrFetchMessage, err)
	}

	id, err := uuid.Parse((*msg).GetId())
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("invalid %T id %s: %v", msg, (*msg).GetId(), err)
		bp.Reader.ProcessFailed()
		return nil, ErrFetchMessage
	}
	if (*msg).WasDeleted() {
		err = bp.Repo.Delete(ctx, &id)
		if err != nil {
			log.Logger.WithContext(ctx).Errorf(
				"unable to delete %T %s: %v", msg, (*msg).GetId(), err)
			bp.Reader.ProcessFailed()
			return nil, ErrProcessingFailed
		}

		log.Logger.WithContext(ctx).Infof("deleted %T %s", msg, (*msg).GetId())
		bp.Reader.ProcessSucceeded(ctx)
		return nil, nil
	}

	err = bp.Repo.Save(ctx, *msg)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf(
			"unable to save %T %s: %v", msg, (*msg).GetId(), err)
		bp.Reader.ProcessFailed()
		return nil, ErrProcessingFailed
	}

	log.Logger.WithContext(ctx).Infof("saved %T %s", msg, (*msg).GetId())
	bp.Reader.ProcessSucceeded(ctx)
	return msg, nil
}

func (bp *DefaultBusProcessor[T]) StopProcessing() {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	if !bp.isProcessing {
		return
	}

	log.Logger.Infof("Stopping bus processor for %s", bp.Reader.GetMessageType())
	bp.isProcessing = false
}

func (bp *DefaultBusProcessor[T]) RegisterListener(listener BusListener[T]) BusListenerHandler {
	listenerHandle := bp.nextListenerHandle
	bp.nextListenerHandle++
	bp.Listeners[listenerHandle] = listener
	return listenerHandle
}

func (bp *DefaultBusProcessor[T]) RemoveListener(listenerHandle BusListenerHandler) {
	delete(bp.Listeners, listenerHandle)
}
