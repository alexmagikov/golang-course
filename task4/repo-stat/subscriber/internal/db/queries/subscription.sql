-- name: CreateSubscription :one
INSERT INTO subscriptions (owner, repo)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteSubscriptions :exec
DELETE FROM subscriptions
WHERE owner = $1 AND repo = $2;

-- name: ListSubscriptions :many
SELECT * FROM subscriptions;

-- name: GetSubscriptionsByRepo :one
SELECT * FROM subscriptions
WHERE owner = $1 AND repo = $2 LIMIT 1;