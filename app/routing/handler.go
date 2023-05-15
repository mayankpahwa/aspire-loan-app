package routing

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	ahttp "github.com/mayankpahwa/aspire-loan-app/internal/resources/http"
	"github.com/mayankpahwa/aspire-loan-app/internal/resources/types"
	"github.com/mayankpahwa/aspire-loan-app/internal/service"
	"github.com/pkg/errors"
)

func CreateUserHandler(service service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ahttp.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp, err := service.CreateUser(r.Context(), req)
		if err != nil {
			WriteError(r, w, err)
			return
		}
		WriteJSON(w, resp)
	})
}

func GetUserLoansHandler(service service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		resp, err := service.GetUserLoans(r.Context(), userID)
		if err != nil {
			WriteError(r, w, err)
			return
		}
		WriteJSON(w, resp)
	})
}

func GetUserLoanByIDHandler(service service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		loanID := chi.URLParam(r, "loanID")
		resp, err := service.GetUserLoanByID(r.Context(), userID, loanID)
		if err != nil {
			WriteError(r, w, err)
			return
		}
		WriteJSON(w, resp)
	})
}

func CreateUserLoanHandler(service service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ahttp.CreateLoanRequest
		userID := chi.URLParam(r, "userID")
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req.UserID = userID
		resp, err := service.CreateUserLoan(r.Context(), req)
		if err != nil {
			WriteError(r, w, err)
			return
		}
		WriteJSON(w, resp)
	})
}

func UpdateUserLoanByIDHandler(service service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(ContextUserIDKey) != "aspire" {
			WriteError(r, w, errors.Wrap(types.ErrForbidden, "only admin can approve a loan"))
			return
		}
		userID := chi.URLParam(r, "userID")
		loanID := chi.URLParam(r, "loanID")
		var req ahttp.UpdateUserLoanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(r, w, errors.Wrap(types.ErrMalformedRequest, err.Error()))
			return
		}
		req.UserID = userID
		req.LoanID = loanID
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			WriteError(r, w, errors.Wrap(types.ErrMalformedRequest, err.Error()))
			return
		}
		resp, err := service.UpdateUserLoanStatus(r.Context(), req)
		if err != nil {
			WriteError(r, w, err)
			return
		}
		WriteJSON(w, resp)
	})
}

func CreateUserLoanRepaymentHandler(service service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ahttp.CreateUserLoanRepaymentRequest
		userID := chi.URLParam(r, "userID")
		loanID := chi.URLParam(r, "loanID")
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(r, w, errors.Wrap(types.ErrMalformedRequest, err.Error()))
			return
		}
		req.UserID = userID
		req.LoanID = loanID
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			WriteError(r, w, errors.Wrap(types.ErrMalformedRequest, err.Error()))
			return
		}
		err := service.CreateUserLoanRepayment(r.Context(), req)
		if err != nil {
			WriteError(r, w, err)
			return
		}
		// WriteJSON(w, resp)
	})
}
