package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/yusadeol/go-gateway-api/internal/domain"
	"github.com/yusadeol/go-gateway-api/internal/dto"
	"github.com/yusadeol/go-gateway-api/internal/service"
	"net/http"
)

type InvoiceHandler struct {
	invoiceService *service.InvoiceService
}

func NewInvoiceHandler(invoiceService *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{invoiceService: invoiceService}
}

func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateInvoiceInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.Account = &dto.AccountInput{
		APIKey: r.Header.Get("X-API-Key"),
	}

	var output *dto.InvoiceOutput
	output, err = h.invoiceService.CreateInvoice(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(&output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *InvoiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	output, err := h.invoiceService.GetByID(id, r.Header.Get("X-API-Key"))
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedAccess) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(&output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *InvoiceHandler) GetByAccountID(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		http.Error(w, "API Key not found", http.StatusUnauthorized)
		return
	}

	output, err := h.invoiceService.ListByAccountAPIKey(apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(&output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
