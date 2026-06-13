package service

import (
	"book-store/internal/domain"
	"book-store/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type externalUsersResponse struct {
	Users []externalUser `json:"users"`
	Total int            `json:"total"`
}

type externalUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CustomerService interface {
	GetAll(ctx context.Context) ([]domain.Customer, error)
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	Create(ctx context.Context, req domain.CreateCustomerRequest) (*domain.Customer, error)
	Update(ctx context.Context, id string, req domain.UpdateCustomerRequest) (*domain.Customer, error)
	Delete(ctx context.Context, id string) error
	Sync(ctx context.Context) (*domain.SyncResult, error)
}

type customerService struct {
	repo       repository.CustomerRepository
	httpClient *http.Client
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{
		repo:       repo,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
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

func (s *customerService) Sync(ctx context.Context) (*domain.SyncResult, error) {
	const (
		baseURL   = "https://dummyjson.com/users"
		pageLimit = 100
	)

	syncID := uuid.NewString()
	result := &domain.SyncResult{}
	skip := 0
	total := 1

	for skip < total {
		page, err := s.fetchUsersPage(ctx, baseURL, skip, pageLimit)
		if err != nil {
			return nil, err
		}
		total = page.Total

		for _, u := range page.Users {
			c := &domain.Customer{
				ID:    uuid.NewString(),
				Name:  u.FirstName + " " + u.LastName,
				Email: u.Email,
			}
			inserted, err := s.repo.Upsert(ctx, c, syncID)
			if err != nil {
				return nil, err
			}
			if inserted {
				result.Created++
			} else {
				result.Updated++
			}
		}

		skip += pageLimit
	}

	deleted, err := s.repo.DeleteStaleSyncedCustomers(ctx, syncID)
	if err != nil {
		return nil, err
	}
	result.Deleted = int(deleted)

	return result, nil
}

func (s *customerService) fetchUsersPage(ctx context.Context, baseURL string, skip, limit int) (*externalUsersResponse, error) {
	url := fmt.Sprintf("%s?limit=%d&skip=%d", baseURL, limit, skip)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	var page externalUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, err
	}
	return &page, nil
}
