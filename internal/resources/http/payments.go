package ahttp

type CreateUserLoanRepaymentRequest struct {
	UserID               string `json:"-"`
	LoanID               string `json:"loan_id"`
	ScheduledRepaymentID string `json:"scheduled_repayment_id"`
	Amount               int    `json:"aount"`
}
