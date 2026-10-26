package httpapi

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yunweigoo/buji/services/api/internal/config"
	"github.com/yunweigoo/buji/services/api/internal/database"
)

type fakeStore struct {
	fabrics  []database.Fabric
	projects []database.Project
}

func (s *fakeStore) ListFabrics(context.Context) ([]database.Fabric, error)   { return s.fabrics, nil }
func (s *fakeStore) ListProjects(context.Context) ([]database.Project, error) { return s.projects, nil }
func (s *fakeStore) CreateFabric(_ context.Context, name, material, pattern, memo string) (database.Fabric, error) {
	item := database.Fabric{ID: "fabric-1", Name: name, Material: material, Pattern: pattern, Memo: memo}
	s.fabrics = append([]database.Fabric{item}, s.fabrics...)
	return item, nil
}
func (s *fakeStore) CreateProject(_ context.Context, name, kind, memo string) (database.Project, error) {
	item := database.Project{ID: "project-1", Name: name, Kind: kind, Memo: memo, Status: "planned"}
	s.projects = append([]database.Project{item}, s.projects...)
	return item, nil
}

func TestHealth(t *testing.T) {
	handler := NewHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), config.Config{Environment: "test"}, nil, nil)
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", responseRecorder.Code)
	}
	var payload map[string]any
	if err := json.Unmarshal(responseRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["requestId"] == "" {
		t.Fatal("expected requestId")
	}
}

func TestCreateFabric(t *testing.T) {
	store := &fakeStore{}
	handler := NewHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), config.Config{Environment: "test"}, nil, store)
	request := httptest.NewRequest(http.MethodPost, "/v1/fabrics", strings.NewReader(`{"name":"蓝色碎花棉布"}`))
	request.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", responseRecorder.Code, responseRecorder.Body.String())
	}
	if len(store.fabrics) != 1 || store.fabrics[0].Name != "蓝色碎花棉布" {
		t.Fatalf("fabric was not created: %#v", store.fabrics)
	}
}

func TestBootstrapReturnsStoredData(t *testing.T) {
	store := &fakeStore{fabrics: []database.Fabric{{ID: "fabric-1", Name: "格纹布"}}}
	handler := NewHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), config.Config{Environment: "test"}, nil, store)
	request := httptest.NewRequest(http.MethodGet, "/v1/bootstrap", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", responseRecorder.Code, responseRecorder.Body.String())
	}
	if !strings.Contains(responseRecorder.Body.String(), "格纹布") {
		t.Fatalf("expected stored fabric in response: %s", responseRecorder.Body.String())
	}
}
