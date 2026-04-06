-- name: CreatePost :one
INSERT INTO posts  (
  title,
  description,
  link,
  guid,
  pub_date,
  source_id
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5, 
  $6
)
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

-- name: UpdatePost :exec
UPDATE posts
  SET title = $1,
    description = $2,
    link = $3,
    guid = $4,
    pub_date = $5,
    source_id = $6,
    status = 'done'
  WHERE guid = $4
  RETURNING *;

-- name: DeduplicatePosts :one
INSERT INTO posts (
  title,
  description,
  link,
  guid,
  pub_date,
  source_id
) 
  VALUES ($1,$2,$3,$4,$5,$6) 
  ON CONFLICT (guid) DO NOTHING
  RETURNING guid;
