package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type SubscribeUserToPlanParams struct {
	UserID int32
	PlanID int32
}

type SubscribeUserToPlanResult struct {
	UserPlan UserPlan
}

func (store *SQLStore) SubscribeUserToPlan(
	ctx context.Context,
	arg SubscribeUserToPlanParams) (SubscribeUserToPlanResult, error) {

	var result SubscribeUserToPlanResult

	err := store.execTx(ctx, func(q *Queries) error {
		// 1. Check and delete existing plan if any
		if err := handleExistingPlan(ctx, q, arg.UserID); err != nil {
			return fmt.Errorf("handling existing plan: %w", err)
		}

		// 2. Create new subscription
		userPlan, err := createNewSubscription(ctx, q, arg)
		if err != nil {
			return fmt.Errorf("creating new subscription: %w", err)
		}

		result.UserPlan = userPlan
		return nil
	})

	return result, err
}

func handleExistingPlan(ctx context.Context, q *Queries, userID int32) error {
	userPlan, err := q.GetOneUserPlan(ctx, pgtype.Int4{
		Int32: userID,
		Valid: true,
	})

	if err == sql.ErrNoRows {
		return nil // No existing plan to handle
	}
	if err != nil {
		return err
	}

	// Delete existing plan
	return q.DeleteUserPlan(ctx, userPlan.UserID)
}

func createNewSubscription(ctx context.Context, q *Queries, arg SubscribeUserToPlanParams) (UserPlan, error) {
	params := InsertUserPlanParams{
		UserID: pgtype.Int4{
			Int32: arg.UserID,
			Valid: true,
		},
		PlanID: pgtype.Int4{
			Int32: arg.PlanID,
			Valid: true,
		},
	}

	return q.InsertUserPlan(ctx, params)
}
