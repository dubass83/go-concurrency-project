-- name: GetOnePlan :one
SELECT * FROM plans
WHERE id = $1 LIMIT 1;

-- name: GetAllPlans :many
SELECT * FROM plans
ORDER by id
LIMIT $1
OFFSET $2;

-- name: UpdatePlan :one
UPDATE plans
SET
  plan_name = COALESCE(sqlc.narg('plan_name'), plan_name),
  plan_amount = COALESCE(sqlc.narg('plan_amount'), plan_amount),
  updated_at = COALESCE(sqlc.narg('updated_at'), updated_at)
WHERE
  id = sqlc.arg('id')
RETURNING *;

-- name: DeletePlan :exec
DELETE FROM plans
WHERE id = $1;
