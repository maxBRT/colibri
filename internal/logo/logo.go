package logo

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"mime"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"www.github.com/maxbrt/colibri/internal/database"
	ps "www.github.com/maxbrt/colibri/internal/pubsub"
)

type Config struct {
	DB            *database.Queries
	S3Client      *s3.Client
	Bucket        string
	PublicBaseURL string
}

type Service struct {
	db            *database.Queries
	client        *s3.Client
	bucket        string
	publicBaseURL string
}

func NewService(cfg Config) (*Service, error) {
	if cfg.DB == nil {
		return nil, fmt.Errorf("database client is required")
	}
	if cfg.S3Client == nil {
		return nil, fmt.Errorf("s3 client is required")
	}
	if strings.TrimSpace(cfg.Bucket) == "" {
		return nil, fmt.Errorf("s3 bucket is required")
	}

	return &Service{
		db:            cfg.DB,
		client:        cfg.S3Client,
		bucket:        cfg.Bucket,
		publicBaseURL: ps.CDNBaseURL,
	}, nil
}

func (s *Service) SaveLogo(ctx context.Context, sourceID string, mimeType string, payload []byte) (database.Logo, error) {
	if len(payload) == 0 {
		return database.Logo{}, fmt.Errorf("empty payload")
	}

	key := s.buildObjectKey(sourceID, mimeType, payload)

	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(payload),
		ContentLength: aws.Int64(int64(len(payload))),
	}
	if strings.TrimSpace(mimeType) != "" {
		input.ContentType = aws.String(mimeType)
	}
	if _, err := s.client.PutObject(ctx, input); err != nil {
		return database.Logo{}, err
	}

	logoURL := s.publicURL(key)

	logo, err := s.db.UpsertLogo(ctx, database.UpsertLogoParams{
		SourceID:  sourceID,
		ObjectKey: key,
		Url:       logoURL,
		MimeType:  nullString(mimeType),
		SizeBytes: nullInt64(int64(len(payload))),
	})
	if err != nil {
		return database.Logo{}, err
	}

	return logo, nil
}

func (s *Service) buildObjectKey(sourceID, mimeType string, payload []byte) string {
	hash := sha256.Sum256(payload)
	hashStr := hex.EncodeToString(hash[:])
	ext := extensionFromMime(mimeType)
	if ext == "" {
		ext = ".bin"
	}
	safeSource := strings.TrimSpace(sourceID)
	if safeSource == "" {
		safeSource = "unknown"
	}

	return path.Join("logos", safeSource, fmt.Sprintf("%s%s", hashStr, ext))
}

func (s *Service) publicURL(key string) string {
	if s.publicBaseURL != "" {
		return fmt.Sprintf("%s/%s", s.publicBaseURL, key)
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
}

func extensionFromMime(mimeType string) string {
	if strings.TrimSpace(mimeType) == "" {
		return ""
	}
	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(exts) == 0 {
		return ""
	}
	return exts[0]
}

func nullString(val string) sql.NullString {
	if strings.TrimSpace(val) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: val, Valid: true}
}

func nullInt64(val int64) sql.NullInt64 {
	if val == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: val, Valid: true}
}
