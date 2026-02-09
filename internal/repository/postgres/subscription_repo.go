package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/always-tired/crud-subscriptions/internal/domain"
	"github.com/always-tired/crud-subscriptions/internal/usecase"
)

const uniqueViolation = "23505"

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s domain.Subscription) (domain.Subscription, error) {
	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	var endDate *time.Time
	var created domain.Subscription
	if err := r.pool.QueryRow(ctx, query,
		s.ID, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate,
	).Scan(
		&created.ID,
		&created.ServiceName,
		&created.Price,
		&created.UserID,
		&created.StartDate,
		&endDate,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return domain.Subscription{}, domain.ErrDuplicate
		}
		return domain.Subscription{}, fmt.Errorf("repo CreateSubscription: %w", err)
	}
	created.EndDate = endDate
	return created, nil
}

func (r *SubscriptionRepository) Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var s domain.Subscription
	var endDate *time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&endDate,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, domain.ErrNotFound
		}
		return domain.Subscription{}, fmt.Errorf("repo GetSubscription: %w", err)
	}
	s.EndDate = endDate
	return s, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, s domain.Subscription) (domain.Subscription, error) {
	query := `
		UPDATE subscriptions
		SET service_name = $2,
			price = $3,
			user_id = $4,
			start_date = $5,
			end_date = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	var endDate *time.Time
	var updated domain.Subscription
	if err := r.pool.QueryRow(ctx, query,
		s.ID, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate,
	).Scan(
		&updated.ID,
		&updated.ServiceName,
		&updated.Price,
		&updated.UserID,
		&updated.StartDate,
		&endDate,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, domain.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return domain.Subscription{}, domain.ErrDuplicate
		}
		return domain.Subscription{}, fmt.Errorf("repo UpdateSubscription: %w", err)
	}
	updated.EndDate = endDate
	return updated, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("repo DeleteSubscription: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *SubscriptionRepository) List(ctx context.Context, filter usecase.ListFilter) ([]domain.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE ($1::uuid IS NULL OR user_id = $1)
		  AND ($2::text IS NULL OR service_name = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	rows, err := r.pool.Query(ctx, query, filter.UserID, filter.ServiceName, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repo ListSubscriptions: %w", err)
	}
	defer rows.Close()

	res := make([]domain.Subscription, 0)
	for rows.Next() {
		var s domain.Subscription
		var endDate *time.Time
		if err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&endDate,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("repo ListSubscriptions: %w", err)
		}
		s.EndDate = endDate
		res = append(res, s)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("repo ListSubscriptions: %w", rows.Err())
	}
	return res, nil
}

func (r *SubscriptionRepository) Summary(ctx context.Context, filter usecase.SummaryFilter) (int64, error) {
	query := `
		WITH months AS (
			SELECT generate_series($1::date, $2::date, interval '1 month') AS m
		)
		SELECT COALESCE(SUM(s.price), 0)
		FROM months m
		JOIN subscriptions s
		  ON s.start_date <= m.m
		 AND (s.end_date IS NULL OR s.end_date >= m.m)
		WHERE ($3::uuid IS NULL OR s.user_id = $3)
		  AND ($4::text IS NULL OR s.service_name = $4)
	`

	var total int64
	if err := r.pool.QueryRow(ctx, query, filter.Start, filter.End, filter.UserID, filter.ServiceName).Scan(&total); err != nil {
		return 0, fmt.Errorf("repo SummarySubscriptions: %w", err)
	}
	return total, nil
}
