package models

import (
	"github.com/google/uuid"
)

type UserLoan struct {
	ID          uuid.UUID `db:"id"`
	UserID      string    `db:"user_id"`
	Amount      int       `db:"amount"`
	Term        int       `db:"term"`
	DateCreated string    `db:"date_created"`
	Status      string    `db:"status"`
}
