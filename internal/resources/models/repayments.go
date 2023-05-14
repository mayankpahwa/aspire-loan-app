package models

import (
	"github.com/google/uuid"
)

type ScheduledRepayment struct {
	ID     uuid.UUID `db:"id"`
	LoanID uuid.UUID `db:"loan_id"`
	Amount int       `db:"amount"`
	Date   string    `db:"date"`
	Status string    `db:"status"`
}

type Repayment struct {
	ID                   uuid.UUID `db:"id"`
	ScheduledRepaymentID uuid.UUID `db:"scheduled_repayment_id"`
	Amount               int       `db:"amount"`
}
