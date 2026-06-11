package handler

import (
	"book-store/internal/domain"
	"book-store/internal/repository"
	"book-store/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// GetAll godoc
// @Summary      List orders
// @Tags         orders
// @Produce      json
// @Success      200  {array}   domain.Order
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /orders [get]
func (h *OrderHandler) GetAll(c *gin.Context) {
	orders, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GetByID godoc
// @Summary      Get an order (with items)
// @Tags         orders
// @Produce      json
// @Param        id   path      string  true  "Order UUID"
// @Success      200  {object}  domain.Order
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /orders/{id} [get]
func (h *OrderHandler) GetByID(c *gin.Context) {
	order, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// GetByCustomer godoc
// @Summary      List orders for a customer
// @Tags         customers
// @Produce      json
// @Param        id   path      string  true  "Customer UUID"
// @Success      200  {array}   domain.Order
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /customers/{id}/orders [get]
func (h *OrderHandler) GetByCustomer(c *gin.Context) {
	orders, err := h.svc.GetByCustomer(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// Create godoc
// @Summary      Create an order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateOrderRequest  true  "Order with items"
// @Success      201   {object}  domain.Order
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

// UpdateStatus godoc
// @Summary      Update order status
// @Tags         orders
// @Accept       json
// @Param        id    path  string                          true  "Order UUID"
// @Param        body  body  domain.UpdateOrderStatusRequest true  "New status"
// @Success      204
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /orders/{id}/status [patch]
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	var req domain.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), c.Param("id"), req); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Delete godoc
// @Summary      Delete an order
// @Tags         orders
// @Param        id   path  string  true  "Order UUID"
// @Success      204
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /orders/{id} [delete]
func (h *OrderHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}