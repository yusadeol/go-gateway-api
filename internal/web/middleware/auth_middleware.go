package middleware

import (
	"errors"
	"github.com/yusadeol/go-gateway-api/internal/domain"
	"github.com/yusadeol/go-gateway-api/internal/service"
	"net/http"
)

type AuthMiddleware struct {
	accountService *service.AccountService
}

func NewAuthMiddleware(accountService *service.AccountService) *AuthMiddleware {
	return &AuthMiddleware{accountService: accountService}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "API Key not found", http.StatusUnauthorized)
			return
		}

		_, err := m.accountService.FindByAPIKey(apiKey)
		if err != nil {
			if errors.Is(err, domain.ErrAccountNotFound) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		next.ServeHTTP(w, r)
	})
}
