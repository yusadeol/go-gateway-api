package service

import (
	"github.com/yusadeol/go-gateway-api/internal/domain"
	"github.com/yusadeol/go-gateway-api/internal/dto"
)

type InvoiceService struct {
	invoiceRepository domain.InvoiceRepository
	accountService    AccountService
}

func NewInvoiceService(invoiceRepository *domain.InvoiceRepository, accountService *AccountService) *InvoiceService {
	return &InvoiceService{
		invoiceRepository: *invoiceRepository,
		accountService:    *accountService,
	}
}

func (s *InvoiceService) CreateInvoice(input *dto.CreateInvoiceInput) (*dto.InvoiceOutput, error) {
	account, err := s.accountService.FindByAPIKey(input.Account.APIKey)
	if err != nil {
		return nil, err
	}

	var invoice *domain.Invoice
	invoice, err = dto.ToInvoice(input, account.ID)

	invoice.Process()

	if invoice.Status == domain.StatusApproved {
		_, err = s.accountService.UpdateBalance(account.APIKey, invoice.Amount)
		if err != nil {
			return nil, err
		}
	}

	err = s.invoiceRepository.Save(invoice)
	if err != nil {
		return nil, err
	}

	return dto.FromInvoice(invoice), nil
}

func (s *InvoiceService) GetByID(id, apiKey string) (*dto.InvoiceOutput, error) {
	invoice, err := s.invoiceRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	var accountOutput *dto.AccountOutput
	accountOutput, err = s.accountService.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	if accountOutput.ID != invoice.AccountID {
		return nil, domain.ErrUnauthorizedAccess
	}

	return dto.FromInvoice(invoice), nil
}

func (s *InvoiceService) ListByAccountID(accountID string) ([]*dto.InvoiceOutput, error) {
	invoices, err := s.invoiceRepository.FindByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	output := make([]*dto.InvoiceOutput, len(invoices))
	for i, invoice := range invoices {
		output[i] = dto.FromInvoice(invoice)
	}

	return output, nil
}

func (s *InvoiceService) ListByAccountAPIKey(apiKey string) ([]*dto.InvoiceOutput, error) {
	account, err := s.accountService.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	return s.ListByAccountID(account.ID)
}
