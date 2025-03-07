package data

import (
	"context"

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
		var err error
		userPlan, err := q.GetOneUserPlan(ctx, pgtype.Int4{
			Int32: arg.UserID,
			Valid: true,
		})
		if err != nil {
			return err
		}
		if userPlan.UserID.Valid {
			err := q.DeleteUserPlan(ctx, pgtype.Int4{
				Int32: arg.UserID,
				Valid: true,
			})
			if err != nil {
				return err
			}
		}
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
		result.UserPlan, err = q.InsertUserPlan(ctx, params)
		return err
	})
	return result, err
}
