package service

import (
	"book-store/internal/domain"
	"book-store/internal/repository"
	"context"

	"github.com/google/uuid"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Create(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error)
	Update(ctx context.Context, id string, req domain.UpdateProductRequest) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAll(ctx context.Context) ([]domain.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *productService) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *productService) Create(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error) {
	p := &domain.Product{
		ID:    uuid.NewString(),
		Name:  req.Name,
		Price: req.Price,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *productService) Update(ctx context.Context, id string, req domain.UpdateProductRequest) (*domain.Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Price > 0 {
		p.Price = req.Price
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *productService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
