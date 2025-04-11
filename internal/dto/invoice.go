package dto

import (
	"github.com/yusadeol/go-gateway-api/internal/domain"
	"time"
)

type Status string

const (
	StatusPending  = Status(domain.StatusPending)
	StatusApproved = Status(domain.StatusApproved)
	StatusRejected = Status(domain.StatusRejected)
)

type CreateInvoiceInput struct {
	APIKey      string
	Description string           `json:"description"`
	PaymentType string           `json:"payment_type"`
	Amount      float64          `json:"amount"`
	Card        *CreditCardInput `json:"card"`
}

type CreditCardInput struct {
	Number      string `json:"number"`
	CVV         string `json:"cvv"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	HolderName  string `json:"holder_name"`
}

type InvoiceOutput struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"account_id"`
	Status         Status    `json:"status"`
	Description    string    `json:"description"`
	PaymentType    string    `json:"payment_type"`
	Amount         float64   `json:"amount"`
	CardLastDigits string    `json:"card_last_digits"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ToInvoice(createInvoiceInput *CreateInvoiceInput, accountID string) (*domain.Invoice, error) {
	return domain.NewInvoice(
		accountID,
		createInvoiceInput.Amount,
		createInvoiceInput.Description,
		createInvoiceInput.PaymentType,
		&domain.CreditCard{
			Number:      createInvoiceInput.Card.Number,
			CVV:         createInvoiceInput.Card.CVV,
			ExpiryMonth: createInvoiceInput.Card.ExpiryMonth,
			ExpiryYear:  createInvoiceInput.Card.ExpiryYear,
			HolderName:  createInvoiceInput.Card.HolderName,
		},
	)
}

func FromInvoice(invoice *domain.Invoice) *InvoiceOutput {
	return &InvoiceOutput{
		invoice.ID,
		invoice.AccountID,
		Status(invoice.Status),
		invoice.Description,
		invoice.PaymentType,
		invoice.Amount,
		invoice.CardLastDigits,
		invoice.CreatedAt,
		invoice.UpdatedAt,
	}
}
