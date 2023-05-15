package mysql

import (
	"context"
	"fmt"

	"github.com/mayankpahwa/aspire-loan-app/internal/resources/models"
)

// GetUserByID fetches a user by ID
func (r Repo) GetUserByID(ctx context.Context, userID string) (models.User, error) {
	var user models.User
	err := r.GetExecutor(ctx).
		QueryRowContext(ctx, "SELECT id, password FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Password)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// CreateUser inserts an entry into the `users` table
func (r Repo) CreateUser(ctx context.Context, user models.User) error {
	result, err := r.GetExecutor(ctx).
		ExecContext(ctx, "INSERT INTO `users` (`id`, `password`) VALUES (?, ?)", user.ID, user.Password)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Finished inserting user. Rows affected: %d\n", rows)
	return nil
}
