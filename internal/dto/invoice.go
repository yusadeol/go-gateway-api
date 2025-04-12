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
	Account     AccountInput
	Description string           `json:"description"`
	PaymentType string           `json:"payment_type"`
	Amount      float64          `json:"amount"`
	Card        *CreditCardInput `json:"card"`
}

type AccountInput struct {
	APIKey string
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

func ToInvoice(input *CreateInvoiceInput, accountID string) (*domain.Invoice, error) {
	return domain.NewInvoice(
		accountID,
		input.Amount,
		input.Description,
		input.PaymentType,
		&domain.CreditCard{
			Number:      input.Card.Number,
			CVV:         input.Card.CVV,
			ExpiryMonth: input.Card.ExpiryMonth,
			ExpiryYear:  input.Card.ExpiryYear,
			HolderName:  input.Card.HolderName,
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
