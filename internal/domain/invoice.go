package domain

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type Invoice struct {
	ID             string
	AccountID      string
	Status         Status
	Description    string
	PaymentType    string
	Amount         float64
	CardLastDigits string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreditCard struct {
	Number      string
	CVV         string
	ExpiryMonth int
	ExpiryYear  int
	HolderName  string
}

func NewInvoice(
	accountID string,
	amount float64,
	description string,
	paymentType string,
	card *CreditCard,
) (*Invoice, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	cardLastDigits := card.Number[len(card.Number)-4:]

	return &Invoice{
		ID:             uuid.New().String(),
		AccountID:      accountID,
		Status:         StatusPending,
		Description:    description,
		PaymentType:    paymentType,
		Amount:         amount,
		CardLastDigits: cardLastDigits,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (i *Invoice) Process() {
	if i.Amount > 10000 {
		return
	}

	randomSource := rand.New(rand.NewSource(time.Now().Unix()))

	if randomSource.Float64() > 0.7 {
		i.Status = StatusRejected
		return
	}

	i.Status = StatusApproved
}

func (i *Invoice) UpdateStatus(status Status) error {
	if status == StatusPending {
		return ErrInvalidStatus
	}

	i.Status = status
	i.UpdatedAt = time.Now()
	return nil
}
