package service

import (
	"testing"
	"time"

	"github.com/mayankpahwa/aspire-loan-app/internal/resources/models"
	"github.com/stretchr/testify/require"
)

func TestShouldChangeLoanStatus(t *testing.T) {
	tests := []struct {
		useCase             string
		scheduledRepayments []models.ScheduledRepayment
		shouldChangeStatus  bool
	}{
		{
			useCase: "when atleast one of the payments is not paid",
			scheduledRepayments: []models.ScheduledRepayment{
				{Status: "PAID"},
				{Status: "PENDING"},
			},
			shouldChangeStatus: false,
		},
		{
			useCase: "when all the payments are paid",
			scheduledRepayments: []models.ScheduledRepayment{
				{Status: "PAID"},
				{Status: "PAID"},
			},
			shouldChangeStatus: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.useCase, func(t *testing.T) {
			shouldChangeStatus := shouldChangeLoanStatus(tc.scheduledRepayments)
			require.Equal(t, tc.shouldChangeStatus, shouldChangeStatus)
		})
	}
}

func TestGetScheduledRepaymentDates(t *testing.T) {
	tests := []struct {
		name                    string
		term                    int
		startDate               string
		scheduledRepaymentDates []string
	}{
		{
			name:                    "term dates in the same month",
			term:                    3,
			startDate:               "2022-02-07",
			scheduledRepaymentDates: []string{"2022-02-14", "2022-02-21", "2022-02-28"},
		},
		{
			name:                    "term dates in different months but same year",
			term:                    4,
			startDate:               "2022-02-08",
			scheduledRepaymentDates: []string{"2022-02-15", "2022-02-22", "2022-03-01", "2022-03-08"},
		},
		{
			name:                    "term dates in different months and different year",
			term:                    5,
			startDate:               "2022-11-30",
			scheduledRepaymentDates: []string{"2022-12-07", "2022-12-14", "2022-12-21", "2022-12-28", "2023-01-04"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			startDate, _ := time.Parse("2006-01-02", tc.startDate)
			scheduledRepaymentDatesActual := getScheduledRepaymentDates(startDate, tc.term)
			require.Equal(t, tc.scheduledRepaymentDates, scheduledRepaymentDatesActual)
		})
	}
}
