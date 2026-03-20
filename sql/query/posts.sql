-- name: CreatePost :one
INSERT INTO posts  (
  title,
  description,
  link,
  guid,
  pub_date,
  source_id
) VALUES ($1,$2,$3,$4,$5, $6)
ON CONFLICT (guid) 
DO UPDATE SET guid = EXCLUDED.guid
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
  WHERE id == $1;

-- name: GetPost :one
SELECT *
  FROM posts
  WHERE id == $1;

-- name: ListPostsForSource :many
SELECT *
  FROM posts
  WHERE LOWER(source_id) = ANY($1::TEXT[]);

-- name: ListPosts :many
SELECT * FROM posts;
