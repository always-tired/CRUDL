package usecase

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionInput struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     *string
}

type ListFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	Limit       int
	Offset      int
}

type SummaryFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	Start       time.Time
	End         time.Time
}
