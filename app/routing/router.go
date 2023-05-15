package routing

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mayankpahwa/aspire-loan-app/app/config"
	"github.com/mayankpahwa/aspire-loan-app/internal/service"
)

// Handler initializes handler for the app
func Handler(conf config.Config, s service.Service) (http.Handler, error) {
	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Post("/", CreateUserHandler(s))
	})
	r.Route("/users/{userID}/loans", func(r chi.Router) {
		r.Use(AuthMiddleware(s.Repo))
		r.Get("/", ValidationUserIdMiddleware(GetUserLoansHandler(s)))
		r.Post("/", ValidationUserIdMiddleware(CreateUserLoanHandler(s)))
		r.Get("/{loanID}", ValidationUserIdMiddleware(GetUserLoanByIDHandler(s)))
		r.Put("/{loanID}", UpdateUserLoanByIDHandler(s))
		r.Post("/{loanID}/payments", (CreateUserLoanRepaymentHandler(s)))
	})
	return r, nil
}
