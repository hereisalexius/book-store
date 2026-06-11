package repository

import (
	"book-store/internal/domain"
	"context"
	"database/sql"
	"errors"
)

type OrderRepository interface {
	FindAll(ctx context.Context) ([]domain.Order, error)
	FindByID(ctx context.Context, id string) (*domain.Order, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]domain.Order, error)
	Create(ctx context.Context, o *domain.Order) error
	UpdateStatus(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error
}

type postgresOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &postgresOrderRepository{db: db}
}

func (r *postgresOrderRepository) FindAll(ctx context.Context) ([]domain.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, customer_id, order_date, status FROM orders ORDER BY order_date DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.OrderDate, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *postgresOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	var o domain.Order
	err := r.db.QueryRowContext(ctx,
		`SELECT id, customer_id, order_date, status FROM orders WHERE id = $1`, id,
	).Scan(&o.ID, &o.CustomerID, &o.OrderDate, &o.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	items, err := r.findItemsByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return &o, nil
}

func (r *postgresOrderRepository) FindByCustomerID(ctx context.Context, customerID string) ([]domain.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, customer_id, order_date, status FROM orders WHERE customer_id = $1 ORDER BY order_date DESC`,
		customerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.OrderDate, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *postgresOrderRepository) Create(ctx context.Context, o *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO orders (id, customer_id, order_date, status) VALUES ($1, $2, $3, $4) RETURNING id`,
		o.ID, o.CustomerID, o.OrderDate, o.Status,
	).Scan(&o.ID)
	if err != nil {
		return err
	}

	for i := range o.Items {
		item := &o.Items[i]
		err = tx.QueryRowContext(ctx,
			`INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			item.ID, o.ID, item.ProductID, item.Quantity, item.Price,
		).Scan(&item.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *postgresOrderRepository) UpdateStatus(ctx context.Context, id, status string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status = $1 WHERE id = $2`, status, id,
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

func (r *postgresOrderRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM orders WHERE id = $1`, id)
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

func (r *postgresOrderRepository) findItemsByOrderID(ctx context.Context, orderID string) ([]domain.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
