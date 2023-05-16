package mysql

import (
	"context"
	"fmt"

	"github.com/mayankpahwa/aspire-loan-app/internal/resources/models"
)

func (r Repo) GetUserLoans(ctx context.Context, userID string) ([]models.UserLoan, error) {
	results, err := r.GetExecutor(ctx).
		QueryContext(ctx, "SELECT id, amount, term, date_created, status FROM loans WHERE user_id = ?", userID)
	if err != nil {
		return []models.UserLoan{}, err
	}
	loans := make([]models.UserLoan, 0)
	for results.Next() {
		var loan models.UserLoan
		if err := results.Scan(&loan.ID, &loan.Amount, &loan.Term, &loan.DateCreated, &loan.Status); err != nil {
			return []models.UserLoan{}, err
		}
		loans = append(loans, loan)
	}
	return loans, nil
}

func (r Repo) GetUserLoanByID(ctx context.Context, userID, loanID string) (models.UserLoan, error) {
	var loan models.UserLoan
	err := r.GetExecutor(ctx).
		QueryRowContext(ctx, "SELECT id, amount, term, date_created, status FROM loans WHERE user_id = ? AND id = ?", userID, loanID).
		Scan(&loan.ID, &loan.Amount, &loan.Term, &loan.DateCreated, &loan.Status)
	if err != nil {
		return models.UserLoan{}, err
	}
	return loan, nil
}

func (r Repo) InsertUserLoan(ctx context.Context, loanToInsert models.UserLoan) error {
	insertUserLoanQuery := "INSERT INTO `loans` (`id`, `user_id`, `amount`, `term`, `date_created`, `status`) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := r.GetExecutor(ctx).
		ExecContext(ctx, insertUserLoanQuery, loanToInsert.ID, loanToInsert.UserID, loanToInsert.Amount, loanToInsert.Term, loanToInsert.DateCreated, loanToInsert.Status)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Finished inserting loan. Rows affected: %d\n", rows)
	return nil
}

func (r Repo) UpdateUserLoanStatus(ctx context.Context, loanID, status string) error {
	updateQuery := "UPDATE `loans` SET `status` = ? WHERE `id` = ?"

	result, err := r.GetExecutor(ctx).ExecContext(ctx, updateQuery, status, loanID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Finished updating loan. Rows affected: %d\n", rows)
	return nil
}
