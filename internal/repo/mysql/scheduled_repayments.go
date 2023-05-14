package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/mayankpahwa/aspire-loan-app/internal/resources/models"
)

func InsertScheduledRepayments(ctx context.Context, tx *sql.Tx, repayments []models.ScheduledRepayment) error {
	placeholders := make([]string, 0)
	args := make([]interface{}, 0)

	for _, repayment := range repayments {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
		args = append(args, repayment.ID, repayment.LoanID, repayment.Amount, repayment.Date, repayment.Status)
	}
	insertRepaymentsQuery := "INSERT INTO `scheduled_repayments` (`id`, `loan_id`, `amount`, `date`, `status`) VALUES %s"
	query := fmt.Sprintf(insertRepaymentsQuery, strings.Join(placeholders, ", "))

	result, err := tx.ExecContext(ctx, query, args...)
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

func UpdateScheduledRepaymentStatus(ctx context.Context, tx *sql.Tx, scheduledRepaymentID, status string) error {
	result, err := tx.
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

func GetScheduledRepaymentsByLoanID(ctx context.Context, loanID string) ([]models.ScheduledRepayment, error) {
	results, err := GetConnection().Query("SELECT id, status FROM scheduled_repayments WHERE loan_id = ?", loanID)
	if err != nil {
		return []models.ScheduledRepayment{}, err
	}
	scheduledRepayments := make([]models.ScheduledRepayment, 0)
	for results.Next() {
		var scheduledRepayment models.ScheduledRepayment
		if err := results.Scan(&scheduledRepayment.ID, &scheduledRepayment.Status); err != nil {
			return []models.ScheduledRepayment{}, err
		}
		scheduledRepayments = append(scheduledRepayments, scheduledRepayment)
	}
	return scheduledRepayments, nil
}