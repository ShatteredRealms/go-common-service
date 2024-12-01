package characterbus

import (
	"context"
	"sync"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
)

type Service interface {
	GetCharacters(ctx context.Context) (*Characters, error)
	GetCharacterById(ctx context.Context, characterId string) (*Character, error)
	DoesOwnCharacter(ctx context.Context, characterId, ownerId string) (bool, error)
	StartProcessingBus(ctx context.Context)
	StopProcessingBus()
}

type service struct {
	repo               Repository
	characterBus       bus.MessageBusReader[Message]
	shouldProcess      bool
	isProcessing       bool
	concurrentErrCount int

	mu sync.Mutex
}

func NewService(
	repo Repository,
	characterBus bus.MessageBusReader[Message],
) Service {
	return &service{
		repo:         repo,
		characterBus: characterBus,
	}
}

// StartProcessingBus implements CharacterService.
func (d *service) StartProcessingBus(ctx context.Context) {
	d.mu.Lock()
	if d.isProcessing {
		d.mu.Unlock()
		return
	}

	d.mu.Lock()
	d.isProcessing = true
	d.shouldProcess = true
	d.mu.Unlock()

	go func() {
		defer func() {
			d.mu.Lock()
			d.isProcessing = false
			d.mu.Unlock()
		}()

		for d.shouldProcess && d.concurrentErrCount < 5 {
			msg, err := d.characterBus.FetchMessage(ctx)
			if err != nil {
				log.Logger.WithContext(ctx).Errorf("unable to fetch character message: %v", err)
				continue
			}

			if d.shouldProcess {
				d.characterBus.ProcessFailed()
				return
			}

			if msg.Deleted {
				_, err = d.repo.DeleteCharacter(ctx, msg.Id)
				if err != nil {
					log.Logger.WithContext(ctx).Errorf(
						"unable to delete character %s: %v", msg.Id, err)
					d.characterBus.ProcessFailed()
					d.concurrentErrCount++
				} else {
					d.characterBus.ProcessSucceeded(ctx)
					d.concurrentErrCount = 0
				}
			} else {
				d.repo.CreateCharacter(ctx, msg.Id, msg.OwnerId)
				if err != nil {
					log.Logger.WithContext(ctx).Errorf(
						"unable to save character %s: %v", msg.Id, err)
					d.characterBus.ProcessFailed()
					d.concurrentErrCount++
				} else {
					d.characterBus.ProcessSucceeded(ctx)
					d.concurrentErrCount = 0
				}
			}
		}

		if d.concurrentErrCount >= 5 {
			d.shouldProcess = false
			log.Logger.WithContext(ctx).Error("too many errors processing character messages")
		}
	}()
}

// StopProcessingBus implements CharacterService.
func (d *service) StopProcessingBus() {
	d.isProcessing = false
}

// GetCharacterById implements CharacterService.
func (d *service) GetCharacterById(ctx context.Context, characterId string) (*Character, error) {
	return d.repo.GetCharacterById(ctx, characterId)
}

// GetCharacters implements CharacterService.
func (d *service) GetCharacters(ctx context.Context) (*Characters, error) {
	return d.repo.GetCharacters(ctx)
}

// DoesOwnCharacter implements Service.
func (d *service) DoesOwnCharacter(ctx context.Context, characterId string, ownerId string) (bool, error) {
	return d.repo.DoesOwnCharacter(ctx, characterId, ownerId)
}
