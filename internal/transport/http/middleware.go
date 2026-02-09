package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/always-tired/crud-subscriptions/internal/domain"
	"github.com/always-tired/crud-subscriptions/internal/usecase"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func requestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)
			log.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.status),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

func recoverer(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic", slog.Any("error", rec))
					writeError(w, http.StatusInternalServerError, "internal error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidArgument):
		writeError(w, http.StatusBadRequest, "invalid argument")
	case errors.Is(err, domain.ErrDuplicate):
		writeError(w, http.StatusConflict, "already exists")
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func domainToResponse(s domain.Subscription) subscriptionResponse {
	var end *string
	if s.EndDate != nil {
		e := usecase.FormatMonthDate(*s.EndDate)
		end = &e
	}

	return subscriptionResponse{
		ID:          s.ID.String(),
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID.String(),
		StartDate:   usecase.FormatMonthDate(s.StartDate),
		EndDate:     end,
		CreatedAt:   s.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
