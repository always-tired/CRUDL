package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/always-tired/crud-subscriptions/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s domain.Subscription) (domain.Subscription, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error)
	Update(ctx context.Context, s domain.Subscription) (domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter ListFilter) ([]domain.Subscription, error)
	Summary(ctx context.Context, filter SummaryFilter) (int64, error)
}

type Service struct {
	repo SubscriptionRepository
	log  *slog.Logger
}

func NewService(repo SubscriptionRepository, log *slog.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) Create(ctx context.Context, input SubscriptionInput) (domain.Subscription, error) {
	sub, err := s.validateInput(input)
	if err != nil {
		return domain.Subscription{}, err
	}
	sub.ID = uuid.New()

	created, err := s.repo.Create(ctx, sub)
	if err != nil {
		s.log.Error("create subscription", "error", err)
		return domain.Subscription{}, err
	}
	return created, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Error("get subscription", "error", err)
		return domain.Subscription{}, err
	}
	return sub, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input SubscriptionInput) (domain.Subscription, error) {
	sub, err := s.validateInput(input)
	if err != nil {
		return domain.Subscription{}, err
	}
	sub.ID = id

	updated, err := s.repo.Update(ctx, sub)
	if err != nil {
		s.log.Error("update subscription", "error", err)
		return domain.Subscription{}, err
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("delete subscription", "error", err)
		return err
	}
	return nil
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]domain.Subscription, error) {
	list, err := s.repo.List(ctx, filter)
	if err != nil {
		s.log.Error("list subscriptions", "error", err)
		return nil, err
	}
	return list, nil
}

func (s *Service) Summary(ctx context.Context, filter SummaryFilter) (int64, error) {
	if filter.Start.IsZero() || filter.End.IsZero() {
		return 0, fmt.Errorf("%w: start and end are required", domain.ErrInvalidArgument)
	}
	if filter.End.Before(filter.Start) {
		return 0, fmt.Errorf("%w: end must be after start", domain.ErrInvalidArgument)
	}

	total, err := s.repo.Summary(ctx, filter)
	if err != nil {
		s.log.Error("summary subscriptions", "error", err)
		return 0, err
	}
	return total, nil
}
