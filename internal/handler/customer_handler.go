package handler

import (
	"book-store/internal/domain"
	"book-store/internal/repository"
	"book-store/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	svc service.CustomerService
}

func NewCustomerHandler(svc service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

// GetAll godoc
// @Summary      List customers
// @Tags         customers
// @Produce      json
// @Success      200  {array}   domain.Customer
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /customers [get]
func (h *CustomerHandler) GetAll(c *gin.Context) {
	customers, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		errResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, customers)
}

// GetByID godoc
// @Summary      Get a customer
// @Tags         customers
// @Produce      json
// @Param        id   path      string  true  "Customer UUID"
// @Success      200  {object}  domain.Customer
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /customers/{id} [get]
func (h *CustomerHandler) GetByID(c *gin.Context) {
	customer, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			errResp(c, http.StatusNotFound, "customer not found")
			return
		}
		errResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, customer)
}

// Create godoc
// @Summary      Create a customer
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateCustomerRequest  true  "Customer data"
// @Success      201   {object}  domain.Customer
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	var req domain.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResp(c, http.StatusBadRequest, err.Error())
		return
	}
	customer, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		errResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, customer)
}

// Update godoc
// @Summary      Update a customer
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id    path      string                        true  "Customer UUID"
// @Param        body  body      domain.UpdateCustomerRequest  true  "Updated fields"
// @Success      200   {object}  domain.Customer
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	var req domain.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResp(c, http.StatusBadRequest, err.Error())
		return
	}
	customer, err := h.svc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			errResp(c, http.StatusNotFound, "customer not found")
			return
		}
		errResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, customer)
}

// Delete godoc
// @Summary      Delete a customer
// @Tags         customers
// @Param        id   path  string  true  "Customer UUID"
// @Success      204
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			errResp(c, http.StatusNotFound, "customer not found")
			return
		}
		errResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}