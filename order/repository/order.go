package repository

import (
	"context"
	"database/sql"
	"distributed-lock/order/postgres"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Order interface {
	Create(ctx context.Context, order postgres.OrderModel) error
	FindByID(ctx context.Context, ID string) (postgres.OrderModel, error)
	Update(ctx context.Context, input OrderUpdate) error
	Delete(ctx context.Context, order postgres.OrderModel) error
}

type OrderUpdate struct {
	InvoiceID string
	Quantity  int
	Value     float64
}

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderDB struct {
	db *sqlx.DB
}

func NewOrderDB(db *sqlx.DB) OrderDB {
	return OrderDB{
		db: db,
	}
}

func (r *OrderDB) Create(ctx context.Context, order postgres.OrderModel) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO orders (
			id,
			user_id,
			invoice_id,
			status,
			item_id,
			quantity,
			value,
			total,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, order.ID, order.UserID, order.InvoiceID, order.Status, order.ItemID,
		order.Quantity, order.Value, order.Total, order.CreatedAt, order.UpdatedAt, order.DeletedAt)
	if err != nil {
		return fmt.Errorf("orderDB.Create inserting order into database failed: %w", err)
	}
	return nil
}

func (r *OrderDB) FindByID(ctx context.Context, ID string) (postgres.OrderModel, error) {
	var order postgres.OrderModel
	err := r.db.GetContext(ctx, &order,
		`
		SELECT id, user_id, invoice_id, status, item_id, quantity, value, total, created_at, updated_at, deleted_at
		FROM orders
		WHERE id = $1
		`, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return postgres.OrderModel{}, ErrOrderNotFound
		}
		return postgres.OrderModel{}, fmt.Errorf("orderDB.FindByID failed: %w", err)
	}
	return order, nil
}

func (r *OrderDB) Update(ctx context.Context, input OrderUpdate) error {
	return nil
}

func (r *OrderDB) Delete(ctx context.Context, order postgres.OrderModel) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE orders SET
		status = $1,
		deleted_at = $2,
		updated_at = $3
		WHERE id = $4
	`, order.Status, order.DeletedAt, order.UpdatedAt, order.ID)
	if err != nil {
		return fmt.Errorf("invoiceDB.Delete failed: %w", err)
	}
	return nil
}
