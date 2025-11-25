package services

import (
	"fmt"
	"log/slog"

	"github.com/QwaQ-dev/servicesSubscription/internal/repository"
	"github.com/QwaQ-dev/servicesSubscription/internal/structures"
	"github.com/QwaQ-dev/servicesSubscription/pkg/sl"
)

type SubscriptionService struct {
	subscriptionRepo *repository.SubscriptionRepo
	log              *slog.Logger
}

func NewSubsriptionService(
	subscriptionRepo *repository.SubscriptionRepo,
	log *slog.Logger,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		log:              log,
	}
}

func (s *SubscriptionService) CreateSub(subscription *structures.Subscription) (int, error) {
	const op = "services.subscriptionService.CreateSub"
	log := s.log.With("op", op)

	id, err := s.subscriptionRepo.InsertSub(subscription)
	if err != nil {
		log.Error("Failed to create subscription", sl.Err(err))
		return 0, fmt.Errorf("%s:%v", op, err)
	}

	log.Info("Subscription created", slog.Int("id", id))

	return id, nil
}

func (s *SubscriptionService) GetAllSubs() ([]structures.Subscription, error) {
	const op = "services.subscriptionService.GetAllSubs"
	log := s.log.With("op", op)

	subscriptions, err := s.subscriptionRepo.SelectAllSubs()
	if err != nil {
		log.Error("Failed to get all subscriptions", sl.Err(err))
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	return subscriptions, nil
}

func (s *SubscriptionService) GetSubById(id int) (structures.Subscription, error) {
	const op = "services.subscriptionService.GetSubById"
	log := s.log.With("op", op)

	subscription, err := s.subscriptionRepo.SelectSubById(id)
	if err != nil {
		log.Error("Failed to get sub by id", sl.Err(err))
		return subscription, fmt.Errorf("%s:%v", op, err)
	}

	return subscription, nil
}

func (s *SubscriptionService) UpdateSub(subscription *structures.Subscription, id int) error {
	const op = "services.subscriptionService.UpdateSub"
	log := s.log.With("op", op)

	err := s.subscriptionRepo.UpdateSub(subscription, id)
	if err != nil {
		log.Error("Failed to update sub", slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *SubscriptionService) DeleteSub(id int) error {
	const op = "services.subscriptionService.DeleteSub"
	log := s.log.With("op", op)

	err := s.subscriptionRepo.DeleteSub(id)
	if err != nil {
		log.Error("Failed to delete sub", slog.Int("id", id), slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SubscriptionService) Counting(data *structures.Counting) (int, error) {
	const op = "services.subscriptionService.Counting"
	log := s.log.With("op", op)

	total, err := s.subscriptionRepo.SelectSum(data)
	if err != nil {
		log.Error("Failed to count sum", sl.Err(err))
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return total, nil
}
