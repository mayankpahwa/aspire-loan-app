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

type Service struct {
	Repo mysql.Repo
}

func NewService(repo mysql.Repo) Service {
	return Service{
		Repo: repo,
	}
}

func (s Service) CreateUser(ctx context.Context, req ahttp.CreateUserRequest) (ahttp.CreateUserResponse, error) {
	_, err := s.Repo.GetUserByID(ctx, req.ID)
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
	if err := s.Repo.CreateUser(ctx, userToInsert); err != nil {
		return ahttp.CreateUserResponse{}, err
	}
	return ahttp.CreateUserResponse{ID: req.ID}, nil
}

// CreateUserLoan inserts a user loan and its scheduled payments
func (s Service) CreateUserLoan(ctx context.Context, req ahttp.CreateLoanRequest) (ahttp.SingleUserLoanResponse, error) {
	var scheduledRepayments []models.ScheduledRepayment
	loanToInsert := getLoanToInsert(req)
	err := s.Repo.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := s.Repo.InsertUserLoan(ctx, loanToInsert); err != nil {
			return err
		}
		scheduledRepayments = getScheduledRepayments(loanToInsert)
		return s.Repo.InsertScheduledRepayments(ctx, scheduledRepayments)
	})
	if err != nil {
		return ahttp.SingleUserLoanResponse{}, nil
	}
	return createInsertUserLoanResponse(loanToInsert, scheduledRepayments), nil
}

// GetUserLoans fetches all loans for a user
func (s Service) GetUserLoans(ctx context.Context, userID string) (ahttp.GetUserLoansResponse, error) {
	result, err := s.Repo.GetUserLoans(ctx, userID)
	if err != nil {
		return ahttp.GetUserLoansResponse{}, errors.Wrap(err, "failed fetching user loans")
	}
	return createGetUserLoansResponse(result), nil
}

// GetUserLoanByID fetches a user loan by loanID
func (s Service) GetUserLoanByID(ctx context.Context, userID, loanID string) (ahttp.SingleUserLoanResponse, error) {
	result, err := s.Repo.GetUserLoanByID(ctx, userID, loanID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ahttp.SingleUserLoanResponse{}, errors.Wrap(types.ErrNoResourceFound, "loan not found")
		}
		return ahttp.SingleUserLoanResponse{}, errors.Wrap(err, "failed fetching user loan by id")
	}
	scheduledRepayments, err := s.Repo.GetScheduledRepaymentsByLoanID(ctx, result.ID.String())
	if err != nil {
		return ahttp.SingleUserLoanResponse{}, errors.Wrap(err, "failed fetching scheduled repayments")
	}
	return createGetUserLoanByIDResponse(result, scheduledRepayments), nil
}

// UpdateUserLoanStatus updates a loan status
func (s Service) UpdateUserLoanStatus(ctx context.Context, req ahttp.UpdateUserLoanRequest) (ahttp.UpdateUserLoanResponse, error) {
	result, err := s.Repo.GetUserLoanByID(ctx, req.UserID, req.LoanID)
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
	if err := s.Repo.UpdateUserLoanStatus(ctx, req.LoanID, types.LoanStatusApproved.ToString()); err != nil {
		return ahttp.UpdateUserLoanResponse{}, err
	}
	return ahttp.UpdateUserLoanResponse{Status: req.Status}, nil
}

// CreateUserLoanRepayment inserts a loan repayment and updates the scheduled repayment status and loan status
func (s Service) CreateUserLoanRepayment(ctx context.Context, req ahttp.CreateUserLoanRepaymentRequest) error {
	loanFromDB, err := s.Repo.GetUserLoanByID(ctx, req.UserID, req.LoanID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Wrap(types.ErrNoResourceFound, "loan not found")
		}
		return err
	}
	if loanFromDB.Status == types.LoanStatusPending.ToString() {
		return errors.Wrap(types.ErrUnprocessableEntity, fmt.Sprintf("cannot create payment for PENDING loan"))
	}
	if loanFromDB.Status == types.LoanStatusPaid.ToString() {
		return errors.Wrap(types.ErrUnprocessableEntity, fmt.Sprintf("Loan already in PAID state"))
	}
	scheduledRepaymentFromDB, err := s.Repo.GetScheduledRepaymentsByID(ctx, req.ScheduledRepaymentID)
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

	txErr := s.Repo.RunInTransaction(ctx, func(ctx context.Context) error {
		repayment := getRepayment(req)
		if err := s.Repo.InsertRepayment(ctx, repayment); err != nil {
			return err
		}

		if err := s.Repo.UpdateScheduledRepaymentStatus(ctx, req.ScheduledRepaymentID, types.ScheduledRepaymentStatusPaid.ToString()); err != nil {
			return err
		}

		scheduledRepayments, err := s.Repo.GetScheduledRepaymentsByLoanID(ctx, req.LoanID)
		if err != nil {
			return err
		}
		if shouldChangeLoanStatus(scheduledRepayments) {
			if err := s.Repo.UpdateUserLoanStatus(ctx, req.LoanID, types.LoanStatusPaid.ToString()); err != nil {
				return err
			}
		}
		return nil
	})
	return txErr
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
	scheduledRepaymentDates := getScheduledRepaymentDates(loanStartDate, loan.Term)
	for _, termDate := range scheduledRepaymentDates {
		scheduledRepayment := models.ScheduledRepayment{
			ID:     uuid.New(),
			LoanID: loan.ID,
			Amount: repaymentInstallation,
			Date:   termDate,
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

func getScheduledRepaymentDates(startDate time.Time, term int) []string {
	repaymentDates := make([]string, 0)
	for i := 1; i <= term; i++ {
		repaymentDate := startDate.AddDate(0, 0, 7*i).Format("2006-01-02")
		repaymentDates = append(repaymentDates, repaymentDate)
	}
	return repaymentDates
}

func shouldChangeLoanStatus(scheduledRepayments []models.ScheduledRepayment) bool {
	for _, scheduledRepayment := range scheduledRepayments {
		if scheduledRepayment.Status != types.ScheduledRepaymentStatusPaid.ToString() {
			return false
		}
	}
	return true
}
