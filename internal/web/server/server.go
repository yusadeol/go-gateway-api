package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/yusadeol/go-gateway-api/internal/service"
	"github.com/yusadeol/go-gateway-api/internal/web/handlers"
	"github.com/yusadeol/go-gateway-api/internal/web/middleware"
	"net/http"
)

type Server struct {
	router         *chi.Mux
	server         *http.Server
	accountService *service.AccountService
	invoiceService *service.InvoiceService
	port           string
}

func NewServer(accountService *service.AccountService, invoiceService *service.InvoiceService, port string) *Server {
	return &Server{
		router:         chi.NewRouter(),
		accountService: accountService,
		invoiceService: invoiceService,
		port:           port,
	}
}

func (s *Server) ConfigureRoutes() {
	authMiddleware := middleware.NewAuthMiddleware(s.accountService)

	accountHandler := handlers.NewAccountHandler(s.accountService)
	invoiceHandler := handlers.NewInvoiceHandler(s.invoiceService)

	s.router.Post("/accounts", accountHandler.Create)

	s.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handle)

		s.router.Get("/accounts", accountHandler.Get)

		s.router.Post("/invoices", invoiceHandler.Create)
		s.router.Get("/invoices/{id}", invoiceHandler.GetByID)

		s.router.Get("/accounts/invoices", invoiceHandler.GetByAccountID)
	})
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	return s.server.ListenAndServe()
}
