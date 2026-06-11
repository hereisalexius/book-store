package repository

import (
	"book-store/internal/domain"
	"context"
	"database/sql"
	"errors"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]domain.Product, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	Create(ctx context.Context, p *domain.Product) error
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id string) error
}

type postgresProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, price FROM products ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *postgresProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	var p domain.Product
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, price FROM products WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Price)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *postgresProductRepository) Create(ctx context.Context, p *domain.Product) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO products (id, name, price) VALUES ($1, $2, $3) RETURNING id`,
		p.ID, p.Name, p.Price,
	).Scan(&p.ID)
}

func (r *postgresProductRepository) Update(ctx context.Context, p *domain.Product) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE products SET name = $1, price = $2 WHERE id = $3`,
		p.Name, p.Price, p.ID,
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

func (r *postgresProductRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
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
