package service

import (
	"book-store/internal/domain"
	"book-store/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderService interface {
	GetAll(ctx context.Context) ([]domain.Order, error)
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	GetByCustomer(ctx context.Context, customerID string) ([]domain.Order, error)
	Create(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id string, req domain.UpdateOrderStatusRequest) error
	Delete(ctx context.Context, id string) error
}

type orderService struct {
	repo     repository.OrderRepository
	products repository.ProductRepository
}

func NewOrderService(repo repository.OrderRepository, products repository.ProductRepository) OrderService {
	return &orderService{repo: repo, products: products}
}

func (s *orderService) GetAll(ctx context.Context) ([]domain.Order, error) {
	return s.repo.FindAll(ctx)
}

func (s *orderService) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *orderService) GetByCustomer(ctx context.Context, customerID string) ([]domain.Order, error) {
	return s.repo.FindByCustomerID(ctx, customerID)
}

func (s *orderService) Create(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error) {
	items := make([]domain.OrderItem, 0, len(req.Items))
	for _, itemReq := range req.Items {
		product, err := s.products.FindByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s: %w", itemReq.ProductID, err)
		}
		items = append(items, domain.OrderItem{
			ID:        uuid.NewString(),
			ProductID: product.ID,
			Quantity:  itemReq.Quantity,
			Price:     product.Price,
		})
	}

	o := &domain.Order{
		ID:         uuid.NewString(),
		CustomerID: req.CustomerID,
		OrderDate:  time.Now(),
		Status:     "pending",
		Items:      items,
	}
	if err := s.repo.Create(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *orderService) UpdateStatus(ctx context.Context, id string, req domain.UpdateOrderStatusRequest) error {
	return s.repo.UpdateStatus(ctx, id, req.Status)
}

func (s *orderService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
