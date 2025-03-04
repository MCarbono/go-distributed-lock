package model

import (
	"database/sql"
	"time"
)

type InvoiceStatus string

const (
	InvoiceStatusCreated InvoiceStatus = "created"
	InvoiceStatusUpdated InvoiceStatus = "updated"
	InvoiceStatusDeleted InvoiceStatus = "deleted"
)

type InvoiceModel struct {
	ID        string        `db:"id"`
	UserID    string        `db:"user_id"`
	Status    InvoiceStatus `db:"status"`
	Total     float64       `db:"total"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	DeletedAt sql.NullTime  `db:"deleted_at"`
}
