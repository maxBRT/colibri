-- +goose Up
CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR NOT NULL,
  description TEXT,
  link TEXT NOT NULL,
  guid TEXT NOT NULL UNIQUE,
  source_id VARCHAR NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  CONSTRAINT fk_sources FOREIGN KEY (source_id)
    REFERENCES sources(id)
);

-- +goose Down
DROP TABLE IF EXISTS posts;
