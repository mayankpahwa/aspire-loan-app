package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mayankpahwa/aspire-loan-app/internal/repo/mysql"
	ahttp "github.com/mayankpahwa/aspire-loan-app/internal/resources/http"
	"github.com/mayankpahwa/aspire-loan-app/internal/resources/models"
	"github.com/mayankpahwa/aspire-loan-app/internal/resources/types"
	"github.com/pkg/errors"
)

func CreateUser(ctx context.Context, req ahttp.CreateUserRequest) (ahttp.CreateUserResponse, error) {
	_, err := mysql.GetUserByID(ctx, req.ID)
	if err == nil {
		return ahttp.CreateUserResponse{}, errors.Wrap(types.ErrUnprocessableEntity, "user already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return ahttp.CreateUserResponse{}, err
	}
	userToInsert := models.User{
		ID:       req.ID,
		Password: req.Password,
	}
	if err := mysql.CreateUser(ctx, userToInsert); err != nil {
		return ahttp.CreateUserResponse{}, err
	}
	return ahttp.CreateUserResponse{ID: req.ID}, nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return ahttp.SingleUserLoanResponse{}, errors.Wrap(types.ErrNoResourceFound, "loan not found")
		}
		return ahttp.SingleUserLoanResponse{}, errors.Wrap(err, "failed fetching user loan by id")
	}
	return createGetUserLoanByIDResponse(result), nil
}

// UpdateUserLoanStatus updates a loan status
func UpdateUserLoanStatus(ctx context.Context, req ahttp.UpdateUserLoanRequest) (ahttp.UpdateUserLoanResponse, error) {
	result, err := mysql.GetUserLoanByID(ctx, req.UserID, req.LoanID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ahttp.UpdateUserLoanResponse{}, errors.Wrap(types.ErrNoResourceFound, "loan not found")
		}
		return ahttp.UpdateUserLoanResponse{}, errors.Wrap(err, "failed fetching user loan by id")
	}
	if result.Status == types.LoanStatusPaid.ToString() {
		return ahttp.UpdateUserLoanResponse{}, errors.Wrap(types.ErrUnprocessableEntity, "cannot approve a paid loan")
	}
	if result.Status == types.LoanStatusApproved.ToString() {
		return ahttp.UpdateUserLoanResponse{Status: req.Status}, nil
	}
	if err := mysql.UpdateUserLoanStatus(ctx, nil, req.LoanID, types.LoanStatusApproved.ToString()); err != nil {
		return ahttp.UpdateUserLoanResponse{}, err
	}
	return ahttp.UpdateUserLoanResponse{Status: req.Status}, nil
}

// CreateUserLoanRepayment inserts a loan repayment and updates the scheduled repayment status and loan status
func CreateUserLoanRepayment(ctx context.Context, req ahttp.CreateUserLoanRepaymentRequest) error {
	loanFromDB, err := mysql.GetUserLoanByID(ctx, req.UserID, req.LoanID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Wrap(types.ErrNoResourceFound, "loan not found")
		}
		return err
	}
	if loanFromDB.Status != types.LoanStatusApproved.ToString() {
		return errors.Wrap(types.ErrUnprocessableEntity, fmt.Sprintf("Loan not in APPROVED status"))
	}
	scheduledRepaymentFromDB, err := mysql.GetScheduledRepaymentsByID(ctx, req.ScheduledRepaymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Wrap(types.ErrNoResourceFound, "scheduled payment not found")
		}
		return err
	}

	if scheduledRepaymentFromDB.Status == types.ScheduledRepaymentStatusPaid.ToString() {
		return errors.Wrap(types.ErrUnprocessableEntity, fmt.Sprintf("Scheduled repayment is already in PAID status"))
	}

	if req.Amount < scheduledRepaymentFromDB.Amount {
		return errors.Wrap(types.ErrUnprocessableEntity, fmt.Sprintf("Repayment amount should be more than scheduled repayment amount"))
	}

	tx, err := mysql.GetConnection().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	repayment := getRepayment(req)
	if err := mysql.InsertRepayment(ctx, tx, repayment); err != nil {
		tx.Rollback()
		return err
	}

	if err := mysql.UpdateScheduledRepaymentStatus(ctx, tx, req.ScheduledRepaymentID, types.ScheduledRepaymentStatusPaid.ToString()); err != nil {
		tx.Rollback()
		return err
	}

	scheduledRepayments, err := mysql.GetScheduledRepaymentsByLoanID(ctx, tx, req.LoanID)
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
	repaymentInstallation := float64(loan.Amount) / float64(loan.Term)
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
