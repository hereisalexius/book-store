// @title          Book Store API
// @version        1.0
// @description    REST API for a book store backed by PostgreSQL.
// @host           localhost:8080
// @BasePath       /api/v1
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Enter "Bearer <token>" — token issued by Microsoft Entra ID.
package main

import (
	_ "book-store/docs"
	"book-store/internal/config"
	"book-store/internal/handler"
	"book-store/internal/middleware"
	"book-store/internal/repository"
	"book-store/internal/server"
	"book-store/internal/service"
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(config.LoadConfig),
		fx.Provide(newDatabase),
		fx.Provide(
			repository.NewCustomerRepository,
			repository.NewProductRepository,
			repository.NewOrderRepository,
		),
		fx.Provide(
			service.NewCustomerService,
			service.NewProductService,
			service.NewOrderService,
		),
		fx.Provide(
			handler.NewCustomerHandler,
			handler.NewProductHandler,
			handler.NewOrderHandler,
		),
		fx.Provide(middleware.NewAuth),
		fx.Provide(server.NewServer),
		fx.Invoke(registerHooks),
	).Run()
}

func newDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func registerHooks(lc fx.Lifecycle, engine *gin.Engine, cfg *config.Config) {
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: engine,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
