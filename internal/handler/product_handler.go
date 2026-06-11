package handler

import (
	"book-store/internal/domain"
	"book-store/internal/repository"
	"book-store/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// GetAll godoc
// @Summary      List products
// @Tags         products
// @Produce      json
// @Success      200  {array}   domain.Product
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		errResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetByID godoc
// @Summary      Get a product
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product UUID"
// @Success      200  {object}  domain.Product
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	product, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			errResp(c, http.StatusNotFound, "product not found")
			return
		}
		errResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, product)
}

// Create godoc
// @Summary      Create a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateProductRequest  true  "Product data"
// @Success      201   {object}  domain.Product
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResp(c, http.StatusBadRequest, err.Error())
		return
	}
	product, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		errResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, product)
}

// Update godoc
// @Summary      Update a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id    path      string                       true  "Product UUID"
// @Param        body  body      domain.UpdateProductRequest  true  "Updated fields"
// @Success      200   {object}  domain.Product
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResp(c, http.StatusBadRequest, err.Error())
		return
	}
	product, err := h.svc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			errResp(c, http.StatusNotFound, "product not found")
			return
		}
		errResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, product)
}

// Delete godoc
// @Summary      Delete a product
// @Tags         products
// @Param        id   path  string  true  "Product UUID"
// @Success      204
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			errResp(c, http.StatusNotFound, "product not found")
			return
		}
		errResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
