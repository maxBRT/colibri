-- name: CreateSource :one
INSERT INTO sources (
  id,
  name,
  url,
  category
) VALUES ( 
  $1,
  $2,
  $3,
  $4 
)
ON CONFLICT (id)
DO UPDATE SET id = EXCLUDED.id 
RETURNING *;

-- name: GetSource :one
SELECT * FROM sources WHERE id == $1;


-- name: DeleteSource :exec
DELETE FROM sources WHERE id == $1;


-- name: ListSources :many
SELECT * FROM sources;

-- name: ListSourcesWithLogo :many
SELECT 
  s.id,
  s.name,
  s.url,
  s.category,
  s.created_at,
  s.updated_at,
  l.url AS logo_url
FROM sources s
LEFT JOIN logos l ON l.source_id = s.id;

-- name: ListCategories :many
SELECT DISTINCT category FROM sources;


-- name: ListSourcesByCategory :many
SELECT * FROM sources WHERE LOWER(category) = ANY($1::text[]);

-- name: ListSourcesByCategoryWithLogo :many
SELECT 
  s.id,
  s.name,
  s.url,
  s.category,
  s.created_at,
  s.updated_at,
  l.url AS logo_url
FROM sources s
LEFT JOIN logos l ON l.source_id = s.id
WHERE LOWER(s.category) = ANY($1::text[]);

-- name: DeduplicateSources :one
INSERT INTO sources (
  id,
  name,
  url,
  category
) 
VALUES ( $1,$2,$3,$4 ) 
ON CONFLICT (id) DO NOTHING
RETURNING id;
