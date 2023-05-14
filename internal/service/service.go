package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mayankpahwa/aspire-loan-app/internal/repo/mysql"
	ahttp "github.com/mayankpahwa/aspire-loan-app/internal/resources/http"
	"github.com/mayankpahwa/aspire-loan-app/internal/resources/models"
	"github.com/mayankpahwa/aspire-loan-app/internal/resources/types"
	"github.com/pkg/errors"
)

func CreateUser(ctx context.Context, req ahttp.CreateUserRequest) {
	fmt.Printf("Creating user: %s\n", req.UserID)
}

// CreateUserLoan inserts a user loan and its scheduled payments
func CreateUserLoan(ctx context.Context, req ahttp.CreateLoanRequest) (ahttp.SingleUserLoanResponse, error) {
	loanToInsert := getLoanToInsert(req)
	tx, err := mysql.GetConnection().BeginTx(ctx, nil)
	if err != nil {
		return ahttp.SingleUserLoanResponse{}, err
	}
	if err := mysql.InsertUserLoan(ctx, tx, loanToInsert); err != nil {
		tx.Rollback()
		return ahttp.SingleUserLoanResponse{}, err
	}
	scheduledRepayments := getScheduledRepayments(loanToInsert)
	if err := mysql.InsertScheduledRepayments(ctx, tx, scheduledRepayments); err != nil {
		tx.Rollback()
		return ahttp.SingleUserLoanResponse{}, err
	}
	tx.Commit()
	return createInsertUserLoanResponse(loanToInsert, scheduledRepayments), nil
}

// GetUserLoans fetches all loans for a user
func GetUserLoans(ctx context.Context, userID string) (ahttp.GetUserLoansResponse, error) {
	result, err := mysql.GetUserLoans(ctx, userID)
	if err != nil {
		return ahttp.GetUserLoansResponse{}, errors.Wrap(err, "failed fetching user loans")
	}
	return createGetUserLoansResponse(result), nil
}

// GetUserLoanByID fetches a user loan by loanID
func GetUserLoanByID(ctx context.Context, userID, loanID string) (ahttp.SingleUserLoanResponse, error) {
	result, err := mysql.GetUserLoanByID(ctx, userID, loanID)
	if err != nil {
		return ahttp.SingleUserLoanResponse{}, errors.Wrap(err, "failed fetching user loan by id")
	}
	return createGetUserLoanByIDResponse(result), nil
}

// CreateUserLoanRepayment inserts a loan repayment and updates the scheduled repayment status and loan status
func CreateUserLoanRepayment(ctx context.Context, req ahttp.CreateUserLoanRepaymentRequest) error {
	loanFromDB, err := mysql.GetUserLoanByID(ctx, req.UserID, req.LoanID)
	if err != nil {
		return err
	}
	if loanFromDB.Status != types.LoanStatusApproved.ToString() {
		return errors.New("Loan not in approved status")
	}
	scheduledRepaymentFromDB, err := mysql.GetUserLoanByID(ctx, req.UserID, req.LoanID)
	if err != nil {
		return err
	}

	if scheduledRepaymentFromDB.Status == types.ScheduledRepaymentStatusPaid.ToString() {
		return errors.New("Scheduled repayment is already paid")
	}
	if scheduledRepaymentFromDB.Amount > req.Amount {
		return err
	}

	tx, err := mysql.GetConnection().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	repayment := getRepayment(req)
	if err := mysql.InsertPayment(ctx, tx, repayment); err != nil {
		tx.Rollback()
		return err
	}

	if err := mysql.UpdateScheduledRepaymentStatus(ctx, tx, req.ScheduledRepaymentID, types.ScheduledRepaymentStatusPaid.ToString()); err != nil {
		tx.Rollback()
		return err
	}

	scheduledRepayments, err := mysql.GetScheduledRepaymentsByLoanID(ctx, req.LoanID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if shouldChangeLoanStatus(scheduledRepayments) {
		if err := mysql.UpdateUserLoanStatus(ctx, tx, req.LoanID, types.LoanStatusPaid.ToString()); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func getLoanToInsert(requestLoan ahttp.CreateLoanRequest) models.UserLoan {
	return models.UserLoan{
		ID:          uuid.New(),
		UserID:      requestLoan.UserID,
		Amount:      requestLoan.Amount,
		Term:        requestLoan.Term,
		DateCreated: time.Now().Format("2006-01-02"),
		Status:      types.LoanStatusPending.ToString(),
	}
}

func getScheduledRepayments(loan models.UserLoan) []models.ScheduledRepayment {
	scheduledRepayments := make([]models.ScheduledRepayment, 0)
	loanStartDate, _ := time.Parse("2006-01-02", loan.DateCreated)
	repaymentInstallation := loan.Amount / loan.Term
	for i := 1; i <= loan.Term; i++ {
		scheduledRepayment := models.ScheduledRepayment{
			ID:     uuid.New(),
			LoanID: loan.ID,
			Amount: repaymentInstallation,
			Date:   loanStartDate.AddDate(0, 0, 7*i).Format("2006-01-02"),
			Status: types.ScheduledRepaymentStatusPending.ToString(),
		}
		scheduledRepayments = append(scheduledRepayments, scheduledRepayment)
	}
	return scheduledRepayments
}

func getRepayment(scheduledRepayment ahttp.CreateUserLoanRepaymentRequest) models.Repayment {
	scheduledRepaymentID, _ := uuid.Parse(scheduledRepayment.ScheduledRepaymentID)
	return models.Repayment{
		ID:                   uuid.New(),
		ScheduledRepaymentID: scheduledRepaymentID,
		Amount:               scheduledRepayment.Amount,
	}
}

func shouldChangeLoanStatus(scheduledRepayments []models.ScheduledRepayment) bool {
	for _, scheduledRepayment := range scheduledRepayments {
		if scheduledRepayment.Status != types.ScheduledRepaymentStatusPaid.ToString() {
			return false
		}
	}
	return true
}
