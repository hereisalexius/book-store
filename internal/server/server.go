package server

import (
	"book-store/internal/config"
	"book-store/internal/handler"
	"book-store/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	CustomerHandler *handler.CustomerHandler
	ProductHandler  *handler.ProductHandler
	OrderHandler    *handler.OrderHandler
	Config          *config.Config
	Auth            *middleware.Auth
}

func NewServer(p Params) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := engine.Group("/api/v1")
	v1.Use(p.Auth.Handler())

	customers := v1.Group("/customers")
	customers.GET("", p.CustomerHandler.GetAll)
	customers.POST("", p.CustomerHandler.Create)
	customers.POST("/sync", p.Auth.RequireRole("Admin"), p.CustomerHandler.Sync)
	customers.GET("/:id", p.CustomerHandler.GetByID)
	customers.PUT("/:id", p.CustomerHandler.Update)
	customers.DELETE("/:id", p.CustomerHandler.Delete)
	customers.GET("/:id/orders", p.OrderHandler.GetByCustomer)

	products := v1.Group("/products")
	products.GET("", p.ProductHandler.GetAll)
	products.POST("", p.ProductHandler.Create)
	products.GET("/:id", p.ProductHandler.GetByID)
	products.PUT("/:id", p.ProductHandler.Update)
	products.DELETE("/:id", p.ProductHandler.Delete)

	orders := v1.Group("/orders")
	orders.GET("", p.OrderHandler.GetAll)
	orders.POST("", p.OrderHandler.Create)
	orders.GET("/:id", p.OrderHandler.GetByID)
	orders.PATCH("/:id/status", p.OrderHandler.UpdateStatus)
	orders.DELETE("/:id", p.OrderHandler.Delete)

	return engine
}
