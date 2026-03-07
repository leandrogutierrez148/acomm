package orders

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lgutierrez148/acomm/internal/http/api"
	"github.com/lgutierrez148/acomm/internal/inbound"
	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/lgutierrez148/acomm/internal/outbound"
)

type OrdersHandler struct {
	repo interfaces.IOrdersRepository
}

func NewOrdersHandler(repo interfaces.IOrdersRepository) *OrdersHandler {
	return &OrdersHandler{repo: repo}
}

func (h *OrdersHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repo.FindAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetOrdersResponse{Orders: mapToOrdersResponse(orders)})
}

func (h *OrdersHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	order, err := h.repo.FindByID(id)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "order not found")
		return
	}
	api.OKResponse(w, outbound.GetOrderResponse{Order: mapToOrderResponse(order)})
}

func (h *OrdersHandler) HandleGetByCustomerEmail(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	if email == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "email is required")
		return
	}

	orders, err := h.repo.FindByCustomerEmail(email)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetOrdersResponse{Orders: mapToOrdersResponse(orders)})
}

func (h *OrdersHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req inbound.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	order := req.ToDomain()
	if err := h.repo.Create(order); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.CreateOrderResponse{Order: mapToOrderResponse(order)})
}

func (h *OrdersHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req inbound.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	order := req.ToDomain()
	order.ID = id
	if err := h.repo.Update(order); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.UpdateOrderResponse{Order: mapToOrderResponse(order)})
}

func (h *OrdersHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, nil)
}

// ─── mappers ────────────────────────────────────────────────────────────────

func mapToOrderResponse(o *models.Order) outbound.Order {
	items := make([]outbound.ItemOrder, len(o.Items))
	for i, it := range o.Items {
		items[i] = outbound.ItemOrder{ItemID: it.ItemID, OrderID: it.OrderID}
	}
	return outbound.Order{
		ID:              o.ID,
		CustomerEmail:   o.CustomerEmail,
		CustomerName:    o.CustomerName,
		CustomerAddress: o.CustomerAddress,
		CustomerPhone:   o.CustomerPhone,
		Status:          o.Status,
		Items:           items,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}
}

func mapToOrdersResponse(orders []models.Order) []outbound.Order {
	resp := make([]outbound.Order, len(orders))
	for i := range orders {
		resp[i] = mapToOrderResponse(&orders[i])
	}
	return resp
}
