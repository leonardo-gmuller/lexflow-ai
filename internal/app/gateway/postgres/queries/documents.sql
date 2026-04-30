-- name: SaveDocument :exec
INSERT INTO documents (
    id,
    case_id,
    file_name,
    path,
    mime_type,
    size_bytes,
    status,
    chunks_count,
    last_error
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
);

-- name: FindDocumentByID :one
SELECT
    id,
    case_id,
    file_name,
    path,
    mime_type,
    size_bytes,
    status,
    chunks_count,
    last_error,
    created_at,
    updated_at
FROM documents
WHERE id = $1;

-- name: ListDocumentsByCaseID :many
SELECT
    id,
    case_id,
    file_name,
    path,
    mime_type,
    size_bytes,
    status,
    chunks_count,
    last_error,
    created_at,
    updated_at
FROM documents
WHERE case_id = $1
ORDER BY created_at DESC;

-- name: UpdateDocumentStatus :exec
UPDATE documents
SET
    status = $2,
    last_error = $3,
    updated_at = now()
WHERE id = $1;

-- name: UpdateDocumentChunksCount :exec
UPDATE documents
SET
    chunks_count = $2,
    updated_at = now()
WHERE id = $1;