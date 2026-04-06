-- name: UpsertLogo :one
INSERT INTO logos (
  source_id,
  object_key,
  url,
  mime_type,
  size_bytes
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
ON CONFLICT (source_id) DO UPDATE SET
  object_key = EXCLUDED.object_key,
  url = EXCLUDED.url,
  mime_type = EXCLUDED.mime_type,
  size_bytes = EXCLUDED.size_bytes,
  updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetLogoBySourceID :one
SELECT *
FROM logos
WHERE source_id = $1;
