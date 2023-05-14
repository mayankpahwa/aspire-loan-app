package types

type (
	LoanStatus               string
	ScheduledRepaymentStatus string
)

var (
	LoanStatusPending  LoanStatus = "PENDING"
	LoanStatusApproved LoanStatus = "APPROVED"
	LoanStatusPaid     LoanStatus = "PAID"

	ScheduledRepaymentStatusPending ScheduledRepaymentStatus = "PENDING"
	ScheduledRepaymentStatusPaid    ScheduledRepaymentStatus = "PAID"
)

func (ls LoanStatus) ToString() string {
	return string(ls)
}

func (srp ScheduledRepaymentStatus) ToString() string {
	return string(srp)
}
