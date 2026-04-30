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
    COUNT(*) OVER()                AS total_count,
    id,
    name,
    description,
    created_at,
    updated_at
FROM cases
WHERE name ILIKE '%' || sqlc.arg(search_term) || '%'
ORDER BY
CASE WHEN sqlc.arg(sort_by) = 'name' AND sqlc.arg(sort_order) = 'asc' THEN name END ASC,
    CASE WHEN sqlc.arg(sort_by) = 'name' AND sqlc.arg(sort_order) = 'desc' THEN name END DESC,

    CASE WHEN sqlc.arg(sort_by) = 'created_at' AND sqlc.arg(sort_order) = 'asc' THEN created_at END ASC,
    CASE WHEN sqlc.arg(sort_by) = 'created_at' AND sqlc.arg(sort_order) = 'desc' THEN created_at END DESC
LIMIT  sqlc.arg(sql_limit)
OFFSET sqlc.arg(sql_offset);