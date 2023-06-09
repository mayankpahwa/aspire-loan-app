package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/mayankpahwa/aspire-loan-app/internal/resources/models"
)

func (r Repo) InsertScheduledRepayments(ctx context.Context, repayments []models.ScheduledRepayment) error {
	placeholders := make([]string, 0)
	args := make([]interface{}, 0)

	for _, repayment := range repayments {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
		args = append(args, repayment.ID, repayment.LoanID, repayment.Amount, repayment.Date, repayment.Status)
	}
	insertRepaymentsQuery := "INSERT INTO `scheduled_repayments` (`id`, `loan_id`, `amount`, `date`, `status`) VALUES %s"
	query := fmt.Sprintf(insertRepaymentsQuery, strings.Join(placeholders, ", "))

	result, err := r.GetExecutor(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Finished inserting scheduled repayments. Rows affected: %d\n", rows)
	return nil
}

func (r Repo) UpdateScheduledRepaymentStatus(ctx context.Context, scheduledRepaymentID, status string) error {
	result, err := r.GetExecutor(ctx).
		ExecContext(ctx, "UPDATE `scheduled_repayments` SET `status` = ? WHERE `id` = ?", status, scheduledRepaymentID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Finished updating `scheduled_repayment`. Rows affected: %d\n", rows)
	return nil
}

// GetScheduledRepaymentsByLoanID fetches all the scheduled repayments for a loanID
func (r Repo) GetScheduledRepaymentsByLoanID(ctx context.Context, loanID string) ([]models.ScheduledRepayment, error) {
	var results *sql.Rows
	var err error
	fetchQuery := "SELECT id, loan_id, amount, date, status FROM scheduled_repayments WHERE loan_id = ?"

	results, err = r.GetExecutor(ctx).QueryContext(ctx, fetchQuery, loanID)
	if err != nil {
		return []models.ScheduledRepayment{}, err
	}
	scheduledRepayments := make([]models.ScheduledRepayment, 0)
	for results.Next() {
		var scheduledRepayment models.ScheduledRepayment
		if err := results.Scan(&scheduledRepayment.ID, &scheduledRepayment.LoanID, &scheduledRepayment.Amount, &scheduledRepayment.Date, &scheduledRepayment.Status); err != nil {
			return []models.ScheduledRepayment{}, err
		}
		scheduledRepayments = append(scheduledRepayments, scheduledRepayment)
	}
	return scheduledRepayments, nil
}

// GetScheduledRepaymentsByID fetches a scheduled repayment by ID
func (r Repo) GetScheduledRepaymentsByID(ctx context.Context, scheduledRepaymentID string) (models.ScheduledRepayment, error) {
	var scheduledRepayment models.ScheduledRepayment
	err := r.GetExecutor(ctx).
		QueryRowContext(ctx, "SELECT id, loan_id, amount, date, status FROM scheduled_repayments WHERE id = ?", scheduledRepaymentID).
		Scan(&scheduledRepayment.ID, &scheduledRepayment.LoanID, &scheduledRepayment.Amount, &scheduledRepayment.Date, &scheduledRepayment.Status)
	if err != nil {
		return models.ScheduledRepayment{}, err
	}
	return scheduledRepayment, nil
}
