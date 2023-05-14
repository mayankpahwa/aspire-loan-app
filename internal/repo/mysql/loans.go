package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mayankpahwa/aspire-loan-app/internal/resources/models"
)

func GetUserLoans(ctx context.Context, userID string) ([]models.UserLoan, error) {
	results, err := GetConnection().Query("SELECT id, amount, term, date_created, status FROM loans WHERE user_id = ?", userID)
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

func GetUserLoanByID(ctx context.Context, userID, loanID string) (models.UserLoan, error) {
	var loan models.UserLoan
	err := GetConnection().
		QueryRow("SELECT id, amount, term, date_created, status FROM loans WHERE user_id = ? AND id = ?", userID, loanID).
		Scan(&loan.ID, &loan.Amount, &loan.Term, &loan.DateCreated, &loan.Status)
	if err != nil {
		return models.UserLoan{}, err
	}
	return loan, nil
}

func InsertUserLoan(ctx context.Context, tx *sql.Tx, loanToInsert models.UserLoan) error {
	result, err := tx.
		ExecContext(ctx, "INSERT INTO `loans` (`id`, `user_id`, `amount`, `term`, `date_created`, `status`) VALUES (?, ?, ?, ?, ?, ?)", loanToInsert.ID, loanToInsert.UserID, loanToInsert.Amount, loanToInsert.Term, loanToInsert.DateCreated, loanToInsert.Status)
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

func UpdateUserLoanStatus(ctx context.Context, tx *sql.Tx, loanID, status string) error {
	var result sql.Result
	var err error
	updateQuery := "UPDATE `loans` SET `status` = ? WHERE `id` = ?"

	if tx == nil {
		result, err = GetConnection().ExecContext(ctx, updateQuery, status, loanID)
	} else {
		result, err = tx.ExecContext(ctx, updateQuery, status, loanID)
	}
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
