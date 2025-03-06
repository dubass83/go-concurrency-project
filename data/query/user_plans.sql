-- name: InsertUserPlan :one
INSERT INTO user_plans (
  user_id, plan_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetOneUserPlan :one
SELECT * FROM user_plans
WHERE id = $1 LIMIT 1;

-- name: GetAllUserPlans :many
SELECT * FROM user_plans
ORDER by id
LIMIT $1
OFFSET $2;

-- name: UpdateUserPlanByUserID :one
UPDATE user_plans
SET
  plan_id = COALESCE(sqlc.narg('plan_id'), plan_id),
  updated_at = COALESCE(sqlc.narg('updated_at'), updated_at)
WHERE
  user_id = sqlc.arg('user_id')
RETURNING *;

-- name: DeleteUserPlanByUserID :exec
DELETE FROM user_plans
WHERE user_id = $1;
