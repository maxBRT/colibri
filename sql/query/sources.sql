-- name: CreateSource :one
INSERT INTO sources (
  id,
  name,
  url,
  category
) VALUES ( $1,$2,$3,$4 )
ON CONFLICT (id)
DO UPDATE SET id = EXCLUDED.id 
RETURNING *;

-- name: GetSource :one
SELECT * FROM sources WHERE id == $1;


-- name: DeleteSource :exec
DELETE FROM sources WHERE id == $1;


-- name: ListSources :many
SELECT * FROM sources;

-- name: ListCategories :many
SELECT DISTINCT category FROM sources;


-- name: ListSourcesByCategory :many
SELECT * FROM sources WHERE LOWER(category) = ANY($1::text[]);