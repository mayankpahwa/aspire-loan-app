package ahttp

type CreateLoanRequest struct {
	UserID string `json:"-"`
	Amount int    `json:"amount"`
	Term   int    `json:"term"`
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
	ID     string `json:"id"`
	LoanID string `json:"loan_id"`
	Amount int    `json:"amount"`
	Date   string `json:"date"`
	Status string `json:"status"`
}
