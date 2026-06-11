package service

import (
	"book-store/internal/domain"
	"book-store/internal/repository"
	"context"

	"github.com/google/uuid"
)

type CustomerService interface {
	GetAll(ctx context.Context) ([]domain.Customer, error)
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	Create(ctx context.Context, req domain.CreateCustomerRequest) (*domain.Customer, error)
	Update(ctx context.Context, id string, req domain.UpdateCustomerRequest) (*domain.Customer, error)
	Delete(ctx context.Context, id string) error
}

type customerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) GetAll(ctx context.Context) ([]domain.Customer, error) {
	return s.repo.FindAll(ctx)
}

func (s *customerService) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *customerService) Create(ctx context.Context, req domain.CreateCustomerRequest) (*domain.Customer, error) {
	c := &domain.Customer{
		ID:    uuid.NewString(),
		Name:  req.Name,
		Email: req.Email,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *customerService) Update(ctx context.Context, id string, req domain.UpdateCustomerRequest) (*domain.Customer, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		c.Name = req.Name
	}
	if req.Email != "" {
		c.Email = req.Email
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *customerService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
