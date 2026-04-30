-- name: SaveDocumentChunk :exec
INSERT INTO document_chunks (
    id,
    document_id,
    case_id,
    content,
    chunk_index,
    token_count,
    embedding
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: ListChunksByDocumentID :many
SELECT
    id,
    document_id,
    case_id,
    content,
    chunk_index,
    token_count,
    embedding,
    created_at
FROM document_chunks
WHERE document_id = $1
ORDER BY chunk_index ASC;

-- name: SearchChunksByEmbedding :many
SELECT
    id,
    document_id,
    case_id,
    content,
    chunk_index,
    token_count,
    created_at,
    embedding <=> $2 AS distance
FROM document_chunks
WHERE case_id = $1
ORDER BY embedding <=> $2
LIMIT $3;