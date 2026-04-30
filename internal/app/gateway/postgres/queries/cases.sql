-- name: CreateCase :exec
INSERT INTO cases (
    id,
    name,
    description
) VALUES (
    $1, $2, $3
);

-- name: FindCaseByID :one
SELECT
    id,
    name,
    description,
    created_at,
    updated_at
FROM cases
WHERE id = $1;

-- name: ListCases :many
SELECT
    id,
    name,
    description,
    created_at,
    updated_at
FROM cases
ORDER BY created_at DESC;