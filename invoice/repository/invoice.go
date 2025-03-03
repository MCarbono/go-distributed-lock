package repository

import (
	"context"
	"database/sql"
	"distributed-lock/invoice/postgres"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrInvoiceNotFound = errors.New("invoice not found")
)

type Invoice interface {
	Create(ctx context.Context, invoice postgres.InvoiceModel) error
	FindByID(ctx context.Context, ID string) (postgres.InvoiceModel, error)
	Update(ctx context.Context, ID string) error
	Delete(ctx context.Context, invoice postgres.InvoiceModel) error
}

type InvoiceDB struct {
	db *sqlx.DB
}

func NewInvoiceDB(db *sqlx.DB) InvoiceDB {
	return InvoiceDB{
		db: db,
	}
}

func (r InvoiceDB) Create(ctx context.Context, invoice postgres.InvoiceModel) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO invoices (
			id,
			user_id,
			status,
			total,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, invoice.ID, invoice.UserID, invoice.Status, invoice.Total, invoice.CreatedAt, invoice.UpdatedAt, invoice.DeletedAt)
	if err != nil {
		return fmt.Errorf("invoiceDB.Create inserting merchant into database failed: %w", err)
	}
	return nil
}

func (r InvoiceDB) FindByID(ctx context.Context, ID string) (postgres.InvoiceModel, error) {
	var invoice postgres.InvoiceModel
	err := r.db.GetContext(ctx, &invoice,
		`
		SELECT id, user_id, status, total, created_at, updated_at, deleted_at
		FROM invoices
		WHERE id = $1
		`, ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return postgres.InvoiceModel{}, ErrInvoiceNotFound
		}
		return postgres.InvoiceModel{}, fmt.Errorf("invoiceDB.FindByID failed: %w", err)
	}
	return invoice, nil
}

func (r InvoiceDB) Update(ctx context.Context, ID string) error {

	return nil
}

func (r InvoiceDB) Delete(ctx context.Context, invoice postgres.InvoiceModel) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE invoices SET
		status = $1,
		deleted_at = $2,
		updated_at = $3
		WHERE id = $4
	`, invoice.Status, invoice.DeletedAt, invoice.UpdatedAt, invoice.ID)
	if err != nil {
		return fmt.Errorf("invoiceDB.Delete failed: %w", err)
	}
	return nil
}
