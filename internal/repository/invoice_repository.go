package repository

import (
	"database/sql"
	"errors"
	"github.com/yusadeol/go-gateway-api/internal/domain"
	"time"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (i *InvoiceRepository) Save(invoice *domain.Invoice) error {
	stmt, err := i.db.Prepare(`
		INSERT INTO invoices (id, account_id, status, description, payment_type, amount, card_last_digits, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		invoice.ID,
		invoice.AccountID,
		invoice.Status,
		invoice.Description,
		invoice.PaymentType,
		invoice.Amount,
		invoice.CardLastDigits,
		invoice.CreatedAt,
		invoice.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (i *InvoiceRepository) FindByID(id string) (*domain.Invoice, error) {
	var invoice domain.Invoice

	err := i.db.QueryRow(
		`SELECT id, account_id, status, description, payment_type, amount, card_last_digits, created_at, updated_at
		FROM invoices WHERE id = ?`,
		id,
	).Scan(
		&invoice.ID,
		&invoice.AccountID,
		&invoice.Status,
		&invoice.Description,
		&invoice.PaymentType,
		&invoice.Amount,
		&invoice.CardLastDigits,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrInvoiceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (i *InvoiceRepository) FindByAccountID(accountID string) ([]*domain.Invoice, error) {
	rows, err := i.db.Query(
		`SELECT id, account_id, status, description, payment_type, amount, card_last_digits, created_at, updated_at
		FROM invoices WHERE account_id = ?`,
		accountID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrInvoiceNotFound
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var invoices []*domain.Invoice
	for rows.Next() {
		var invoice domain.Invoice
		err := rows.Scan(
			&invoice.ID,
			&invoice.AccountID,
			&invoice.Status,
			&invoice.Description,
			&invoice.PaymentType,
			&invoice.Amount,
			&invoice.CardLastDigits,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		invoices = append(invoices, &invoice)
	}

	return invoices, nil
}

func (i *InvoiceRepository) UpdateStatus(invoice *domain.Invoice) error {
	var res sql.Result
	res, err := i.db.Exec(
		`UPDATE invoices SET status = ?, updated_at = ? WHERE id = ?`,
		invoice.Status,
		time.Now(),
		invoice.ID,
	)

	if err != nil {
		return err
	}

	var affectedRows int64
	affectedRows, err = res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return domain.ErrInvoiceNotFound
	}

	return nil
}
