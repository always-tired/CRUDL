package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/always-tired/crud-subscriptions/internal/usecase"
)

type Handler struct {
	service *usecase.Service
	log     *slog.Logger
}

func NewHandler(service *usecase.Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(requestLogger(h.log))
	r.Use(recoverer(h.log))

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.createSubscription)
		r.Get("/", h.listSubscriptions)
		r.Get("/summary", h.summary)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.getSubscription)
			r.Put("/", h.updateSubscription)
			r.Delete("/", h.deleteSubscription)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return r
}

// @Summary Create subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body subscriptionRequest true "subscription"
// @Success 201 {object} subscriptionResponse
// @Failure 400 {object} errorResponse
// @Router /subscriptions [post]
func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	var req subscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	input := usecase.SubscriptionInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	created, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, domainToResponse(created))
}

// @Summary Get subscription by id
// @Tags subscriptions
// @Produce json
// @Param id path string true "subscription id" format(uuid)
// @Success 200 {object} subscriptionResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /subscriptions/{id} [get]
func (h *Handler) getSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	sub, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, domainToResponse(sub))
}

// @Summary Update subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "subscription id" format(uuid)
// @Param subscription body subscriptionRequest true "subscription"
// @Success 200 {object} subscriptionResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /subscriptions/{id} [put]
func (h *Handler) updateSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req subscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	input := usecase.SubscriptionInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	updated, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, domainToResponse(updated))
}

// @Summary Delete subscription
// @Tags subscriptions
// @Param id path string true "subscription id" format(uuid)
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /subscriptions/{id} [delete]
func (h *Handler) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary List subscriptions
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "user id" format(uuid)
// @Param service_name query string false "service name"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {array} subscriptionResponse
// @Router /subscriptions [get]
func (h *Handler) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	var filter usecase.ListFilter

	if v := r.URL.Query().Get("user_id"); v != "" {
		uid, err := uuid.Parse(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		filter.UserID = &uid
	}
	if v := r.URL.Query().Get("service_name"); v != "" {
		filter.ServiceName = &v
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Offset = n
		}
	}

	list, err := h.service.List(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	resp := make([]subscriptionResponse, 0, len(list))
	for _, s := range list {
		resp = append(resp, domainToResponse(s))
	}

	writeJSON(w, http.StatusOK, resp)
}

// @Summary Get total cost for period
// @Tags subscriptions
// @Produce json
// @Param start query string true "start month" example(07-2025)
// @Param end query string true "end month" example(12-2025)
// @Param user_id query string false "user id" format(uuid)
// @Param service_name query string false "service name"
// @Success 200 {object} map[string]int64
// @Failure 400 {object} errorResponse
// @Router /subscriptions/summary [get]
func (h *Handler) summary(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	if startStr == "" || endStr == "" {
		writeError(w, http.StatusBadRequest, "start and end are required")
		return
	}

	start, err := usecase.ParseMonthDate(startStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	end, err := usecase.ParseMonthDate(endStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	filter := usecase.SummaryFilter{Start: start, End: end}

	if v := r.URL.Query().Get("user_id"); v != "" {
		uid, err := uuid.Parse(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		filter.UserID = &uid
	}
	if v := r.URL.Query().Get("service_name"); v != "" {
		filter.ServiceName = &v
	}

	total, err := h.service.Summary(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"total": total,
	})
}
