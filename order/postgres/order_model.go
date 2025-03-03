package postgres

import (
	"database/sql"
	"time"
)

type OrderStatus string

const (
	OrderStatusCreated OrderStatus = "created"
	OrderStatusUpdated OrderStatus = "updated"
	OrderStatusDeleted OrderStatus = "deleted"
)

type OrderModel struct {
	ID        string       `db:"id"`
	UserID    string       `db:"user_id"`
	InvoiceID string       `db:"invoice_id"`
	Status    OrderStatus  `db:"status"`
	ItemID    string       `db:"item_id"`
	Quantity  int          `db:"quantity"`
	Value     float64      `db:"value"`
	Total     float64      `db:"total"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
