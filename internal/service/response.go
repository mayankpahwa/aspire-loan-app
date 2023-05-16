package service

import (
	ahttp "github.com/mayankpahwa/aspire-loan-app/internal/resources/http"
	"github.com/mayankpahwa/aspire-loan-app/internal/resources/models"
)

func createGetUserLoansResponse(response []models.UserLoan) ahttp.GetUserLoansResponse {
	var resp ahttp.GetUserLoansResponse
	userLoans := make([]ahttp.SingleUserLoanResponse, 0)
	for _, loan := range response {
		singleUserLoan := ahttp.SingleUserLoanResponse{
			ID:          loan.ID.String(),
			Amount:      loan.Amount,
			Term:        loan.Term,
			DateCreated: loan.DateCreated,
			Status:      loan.Status,
		}
		userLoans = append(userLoans, singleUserLoan)
	}
	resp.UserLoans = userLoans
	return resp
}

func createGetUserLoanByIDResponse(response models.UserLoan, scheduledRepayments []models.ScheduledRepayment) ahttp.SingleUserLoanResponse {
	scheduledRepaymentsResp := make([]ahttp.SingleScheduledRepayment, 0)
	for _, scheduledRepayment := range scheduledRepayments {
		singleScheduledRepayment := ahttp.SingleScheduledRepayment{
			ID:     scheduledRepayment.ID.String(),
			LoanID: scheduledRepayment.LoanID.String(),
			Amount: scheduledRepayment.Amount,
			Date:   scheduledRepayment.Date,
			Status: scheduledRepayment.Status,
		}
		scheduledRepaymentsResp = append(scheduledRepaymentsResp, singleScheduledRepayment)
	}
	return ahttp.SingleUserLoanResponse{
		ID:                  response.ID.String(),
		Amount:              response.Amount,
		Term:                response.Term,
		DateCreated:         response.DateCreated,
		Status:              response.Status,
		ScheduledRepayments: scheduledRepaymentsResp,
	}
}

func createInsertUserLoanResponse(loan models.UserLoan, scheduledRepayments []models.ScheduledRepayment) ahttp.SingleUserLoanResponse {
	scheduledRepaymentsResponse := make([]ahttp.SingleScheduledRepayment, 0)
	for _, scheduledRepayment := range scheduledRepayments {
		scheduledRepaymentsResponse = append(scheduledRepaymentsResponse, ahttp.SingleScheduledRepayment{
			ID:     scheduledRepayment.ID.String(),
			LoanID: scheduledRepayment.LoanID.String(),
			Amount: scheduledRepayment.Amount,
			Date:   scheduledRepayment.Date,
			Status: scheduledRepayment.Status,
		})
	}
	return ahttp.SingleUserLoanResponse{
		ID:                  loan.ID.String(),
		Amount:              loan.Amount,
		Term:                loan.Term,
		DateCreated:         loan.DateCreated,
		Status:              loan.Status,
		ScheduledRepayments: scheduledRepaymentsResponse,
	}
}
