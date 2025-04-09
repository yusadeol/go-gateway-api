package dto

import (
	"github.com/yusadeol/go-gateway-api/internal/domain"
	"time"
)

type CreateAccountInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	APIKey    string    `json:"api_key"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToAccount(input *CreateAccountInput) *domain.Account {
	return domain.NewAccount(input.Name, input.Email)
}

func FromAccount(account *domain.Account) *AccountOutput {
	return &AccountOutput{
		ID:        account.ID,
		Name:      account.Name,
		Email:     account.Email,
		APIKey:    account.APIKey,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}
