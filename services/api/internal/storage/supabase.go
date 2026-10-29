package storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const Bucket = "fabric-images"

type Supabase struct {
	baseURL string
	key     string
	client  *http.Client
}

func NewSupabase(baseURL, key string) *Supabase {
	return &Supabase{
		baseURL: strings.TrimRight(baseURL, "/"),
		key:     key,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *Supabase) Ready() bool { return s != nil && s.baseURL != "" && s.key != "" }

func (s *Supabase) EnsureBucket(ctx context.Context) error {
	if !s.Ready() {
		return errors.New("supabase storage is not configured")
	}
	body, _ := json.Marshal(map[string]any{"id": Bucket, "name": Bucket, "public": true, "file_size_limit": 20 << 20})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/storage/v1/bucket", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	s.authorize(req)
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("create storage bucket: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusConflict {
		return nil
	}
	message, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	if resp.StatusCode == http.StatusBadRequest && bytes.Contains(bytes.ToLower(message), []byte("already exists")) {
		return nil
	}
	return fmt.Errorf("create storage bucket: status %d: %s", resp.StatusCode, strings.TrimSpace(string(message)))
}

func (s *Supabase) Upload(ctx context.Context, ownerID, extension, contentType string, data []byte) (string, error) {
	random := make([]byte, 16)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	key := ownerID + "/" + hex.EncodeToString(random) + extension
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/storage/v1/object/"+Bucket+"/"+key, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "false")
	s.authorize(req)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload object: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return "", fmt.Errorf("upload object: status %d: %s", resp.StatusCode, strings.TrimSpace(string(message)))
	}
	return key, nil
}

func (s *Supabase) PublicURL(key string) string {
	if key == "" {
		return ""
	}
	return s.baseURL + "/storage/v1/object/public/" + Bucket + "/" + key
}

func (s *Supabase) authorize(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+s.key)
	req.Header.Set("apikey", s.key)
}
