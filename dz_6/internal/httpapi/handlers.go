package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"example.com/lesson8/internal/model"
	"example.com/lesson8/internal/repository"
	"example.com/lesson8/internal/service"
)

type Orders interface {
	Create(customer string) (model.Order, error)
	ByID(id int) (model.Order, error)
}

type API struct {
	orders Orders
	logger *slog.Logger
}

func New(orders Orders, logger *slog.Logger) *API { return &API{orders: orders, logger: logger} }

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	a.writeJSON(w, r, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Customer string `json:"customer"`
	}

	// Здесь проверяется формат. Где должно жить правило о пустом имени?
	if err := decodeJSON(w, r, &in); err != nil {
		a.writeError(w, r, err.status, err.message)
		return
	}
	order, err := a.orders.Create(in.Customer)
	if errors.Is(err, service.ErrEmptyCustomer) {
		a.writeError(w, r, http.StatusBadRequest, "customer is required")
		return
	}
	if err != nil {
		a.internalError(w, r, err)
		return
	}
	w.Header().Set("Location", "/orders/"+strconv.Itoa(order.ID))
	a.writeJSON(w, r, http.StatusCreated, order)
}

func (a *API) getOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		a.writeError(w, r, http.StatusBadRequest, "id must be a positive integer")
		return
	}
	order, err := a.orders.ByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		a.writeError(w, r, http.StatusNotFound, "order not found")
		return
	}
	if err != nil {
		a.internalError(w, r, err)
		return
	}
	a.writeJSON(w, r, http.StatusOK, order)
}

func (a *API) internalError(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.ErrorContext(r.Context(), "order operation failed", "error", err, "request_id", RequestID(r.Context()))
	a.writeError(w, r, http.StatusInternalServerError, "internal server error")
}

func (a *API) writeError(w http.ResponseWriter, r *http.Request, status int, message string) {
	a.writeJSON(w, r, status, struct {
		Error string `json:"error"`
	}{message})
}

func (a *API) writeJSON(w http.ResponseWriter, r *http.Request, status int, value any) {
	// Что произойдёт, если сериализация завершится ошибкой после WriteHeader?
	data, err := json.Marshal(value)
	if err != nil {
		a.logger.ErrorContext(r.Context(), "encode response failed", "error", err, "request_id", RequestID(r.Context()))
		status, data = http.StatusInternalServerError, []byte(`{"error":"internal server error"}`)
		w.Header().Del("Location")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(append(data, '\n')); err != nil {
		// Статус уже отправлен. Второй ответ здесь записывать нельзя.
		a.logger.ErrorContext(r.Context(), "write response failed", "error", err, "request_id", RequestID(r.Context()))
	}
}
