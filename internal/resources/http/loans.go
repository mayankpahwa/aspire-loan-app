package ahttp

type CreateLoanRequest struct {
	UserID string `json:"-"`
	Amount int    `json:"amount" validate:"required,gt=0"`
	Term   int    `json:"term" validate:"required,gt=0"`
}

type CreateLoanResponse struct {
	ID          string `json:"id"`
	Amount      int    `json:"amount"`
	Term        int    `json:"term"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

type GetUserLoansResponse struct {
	UserLoans []SingleUserLoanResponse `json:"loans"`
}

type SingleUserLoanResponse struct {
	ID                  string                     `json:"id"`
	Amount              int                        `json:"amount"`
	Term                int                        `json:"term"`
	DateCreated         string                     `json:"date_created"`
	Status              string                     `json:"status"`
	ScheduledRepayments []SingleScheduledRepayment `json:"scheduled_repayments,omitempty"`
}

type SingleScheduledRepayment struct {
	ID     string  `json:"id"`
	LoanID string  `json:"loan_id"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
	Status string  `json:"status"`
}

type UpdateUserLoanRequest struct {
	UserID string `json:"-"`
	LoanID string `json:"-"`
	Status string `json:"status" validate:"required,oneof=APPROVED"`
}

type UpdateUserLoanResponse struct {
	Status string `json:"status"`
}
