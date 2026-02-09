package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/always-tired/crud-subscriptions/internal/domain"
)

func (s *Service) validateInput(input SubscriptionInput) (domain.Subscription, error) {
	name := strings.TrimSpace(input.ServiceName)
	if len(name) < 3 {
		return domain.Subscription{}, fmt.Errorf("%w: service_name must be at least 3 characters", domain.ErrInvalidArgument)
	}
	if input.Price <= 0 {
		return domain.Subscription{}, fmt.Errorf("%w: price must be positive integer", domain.ErrInvalidArgument)
	}

	uid, err := uuid.Parse(input.UserID)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("%w: invalid user_id", domain.ErrInvalidArgument)
	}

	start, err := ParseMonthDate(input.StartDate)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("%w: %s", domain.ErrInvalidArgument, err.Error())
	}

	var end *time.Time
	if input.EndDate != nil && strings.TrimSpace(*input.EndDate) != "" {
		t, err := ParseMonthDate(*input.EndDate)
		if err != nil {
			return domain.Subscription{}, fmt.Errorf("%w: %s", domain.ErrInvalidArgument, err.Error())
		}
		end = &t
	}
	if end != nil && end.Before(start) {
		return domain.Subscription{}, fmt.Errorf("%w: end_date must be after start_date", domain.ErrInvalidArgument)
	}

	sub := domain.Subscription{
		ServiceName: name,
		Price:       input.Price,
		UserID:      uid,
		StartDate:   start,
		EndDate:     end,
	}
	if err := s.validateDomain(sub); err != nil {
		return domain.Subscription{}, err
	}
	return sub, nil
}

func (s *Service) validateDomain(sub domain.Subscription) error {
	if sub.UserID == uuid.Nil {
		return fmt.Errorf("%w: user_id is required", domain.ErrInvalidArgument)
	}
	if sub.StartDate.IsZero() {
		return fmt.Errorf("%w: start_date is required", domain.ErrInvalidArgument)
	}
	return nil
}
