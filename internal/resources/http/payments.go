package ahttp

// CreateUserLoanRepaymentRequest represents the request to create a loan repayment
type CreateUserLoanRepaymentRequest struct {
	UserID               string  `json:"-"`
	LoanID               string  `json:"-"`
	ScheduledRepaymentID string  `json:"scheduled_repayment_id" validate:"required"`
	Amount               float64 `json:"amount"`
}
