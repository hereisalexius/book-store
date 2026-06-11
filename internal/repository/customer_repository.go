package repository

import (
	"book-store/internal/domain"
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("not found")

type CustomerRepository interface {
	FindAll(ctx context.Context) ([]domain.Customer, error)
	FindByID(ctx context.Context, id string) (*domain.Customer, error)
	Create(ctx context.Context, c *domain.Customer) error
	Update(ctx context.Context, c *domain.Customer) error
	Delete(ctx context.Context, id string) error
}

type postgresCustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &postgresCustomerRepository{db: db}
}

func (r *postgresCustomerRepository) FindAll(ctx context.Context) ([]domain.Customer, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, email FROM customers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []domain.Customer
	for rows.Next() {
		var c domain.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, rows.Err()
}

func (r *postgresCustomerRepository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	var c domain.Customer
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email FROM customers WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *postgresCustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO customers (id, name, email) VALUES ($1, $2, $3) RETURNING id`,
		c.ID, c.Name, c.Email,
	).Scan(&c.ID)
}

func (r *postgresCustomerRepository) Update(ctx context.Context, c *domain.Customer) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE customers SET name = $1, email = $2 WHERE id = $3`,
		c.Name, c.Email, c.ID,
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresCustomerRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM customers WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
