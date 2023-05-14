package routing

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mayankpahwa/aspire-loan-app/app/config"
)

// Handler initializes handler for the app
func Handler(conf config.Config) (http.Handler, error) {
	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Post("/", CreateUserHandler())
	})
	r.Route("/users/{userID}/loans", func(r chi.Router) {
		r.Use(AuthMiddleware())
		r.Get("/", ValidationUserIdMiddleware(GetUserLoansHandler()))
		r.Post("/", ValidationUserIdMiddleware(CreateUserLoanHandler()))
		r.Get("/{loanID}", ValidationUserIdMiddleware(GetUserLoanByIDHandler()))
		r.Put("/{loanID}", UpdateUserLoanByIDHandler())
		r.Post("/{loanID}/payments", (CreateUserLoanRepaymentHandler()))
	})
	return r, nil
}
