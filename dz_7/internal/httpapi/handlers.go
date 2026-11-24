package httpapi

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"example.com/lesson9practice/internal/model"
	"example.com/lesson9practice/internal/repository"
	"example.com/lesson9practice/internal/service"
)

type Orders interface {
	Create(string) (model.Order, error)
	ByID(int) (model.Order, error)
	List(limit, offset int) ([]model.Order, error)
}

type API struct {
	orders Orders
	logger *slog.Logger
}

func New(orders Orders, logger *slog.Logger) *API { return &API{orders: orders, logger: logger} }

// @Summary Проверка статуса сервиса
// @Tags status
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /status [get]
func (a *API) status(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, StatusResponse{Status: "ok"})
}

// @Summary Создать заказ
// @Tags orders
// @Accept json
// @Produce json
// @Param request body CreateOrderRequest true "Данные заказа"
// @Success 201 {object} OrderResponse
// @Header 201 {string} Location "Путь к созданному заказу"
// @Failure 400 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 415 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [post]
func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	var in *CreateOrderRequest
	if err := decodeJSON(w, r, &in); err != nil {
		writeError(w, err.status, err.message, err.code)
		return
	}
	if in == nil {
		writeError(w, 400, "body must be a JSON object", "INVALID_JSON")
		return
	}
	order, err := a.orders.Create(in.Customer)
	if errors.Is(err, service.ErrEmptyCustomer) {
		writeError(w, 400, "customer is required", "VALIDATION_ERROR")
		return
	}
	if err != nil {
		a.internalError(w, r, err)
		return
	}
	w.Header().Set("Location", "/api/v1/orders/"+strconv.Itoa(order.ID))
	writeJSON(w, 201, orderResponse(order))
}

// @Summary Получить заказ по ID
// @Tags orders
// @Produce json
// @Param id path int true "ID заказа" minimum(1)
// @Success 200 {object} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id} [get]
func (a *API) getOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeError(w, 400, "id must be a positive integer", "INVALID_ID")
		return
	}
	order, err := a.orders.ByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, 404, "order not found", "ORDER_NOT_FOUND")
		return
	}
	if err != nil {
		a.internalError(w, r, err)
		return
	}
	writeJSON(w, 200, orderResponse(order))
}

// @Summary Получить список заказов
// @Tags orders
// @Produce json
// @Param limit query int false "Количество заказов" default(20) minimum(1) maximum(100)
// @Param offset query int false "Смещение" default(0) minimum(0)
// @Success 200 {array} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [get]
func (a *API) listOrders(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := parseListQuery(r.URL.RawQuery)
	if !ok {
		writeError(w, 400, "invalid query parameters", "INVALID_QUERY")
		return
	}
	orders, err := a.orders.List(limit, offset)
	if err != nil {
		a.internalError(w, r, err)
		return
	}
	res := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		res = append(res, orderResponse(o))
	}
	writeJSON(w, 200, res)
}

func parseListQuery(rawQuery string) (int, int, bool) {
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return 0, 0, false
	}
	limit := 20
	offset := 0
	for key, list := range values {
		if len(list) != 1 {
			return 0, 0, false
		}
		switch key {
		case "limit":
			v, ok := parseDigits(list[0])
			if !ok || v < 1 || v > 100 {
				return 0, 0, false
			}
			limit = v
		case "offset":
			v, ok := parseDigits(list[0])
			if !ok {
				return 0, 0, false
			}
			offset = v
		default:
			return 0, 0, false
		}
	}
	return limit, offset, true
}

func parseDigits(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, false
		}
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return v, true
}

func (a *API) internalError(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.ErrorContext(r.Context(), "order operation failed", "error", err)
	writeError(w, 500, "internal server error", "INTERNAL_ERROR")
}
