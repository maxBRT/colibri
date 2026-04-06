-- +goose Up
CREATE TABLE logos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  source_id VARCHAR NOT NULL UNIQUE,
  object_key TEXT NOT NULL,
  url TEXT NOT NULL,
  mime_type TEXT,
  size_bytes BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  CONSTRAINT fk_logos_source FOREIGN KEY (source_id)
    REFERENCES sources(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS logos;
