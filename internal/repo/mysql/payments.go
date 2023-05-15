package mysql

import (
	"context"
	"fmt"

	"github.com/mayankpahwa/aspire-loan-app/internal/resources/models"
)

func (r Repo) InsertRepayment(ctx context.Context, repayments models.Repayment) error {
	insertPaymentQuery := "INSERT INTO `repayments` (`id`, `scheduled_repayments_id`, `amount`) VALUES (?, ?, ?)"

	result, err := r.GetExecutor(ctx).ExecContext(ctx, insertPaymentQuery, repayments.ID, repayments.ScheduledRepaymentID, repayments.Amount)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Finished inserting payments. Rows affected: %d\n", rows)
	return nil
}
