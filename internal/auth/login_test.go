package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin_Success_ReturnsKey(t *testing.T) {
	companiesCalled := false
	createKeyCalled := false

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api-v2/auth/companies":
			companiesCalled = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"content": []map[string]interface{}{
					{"id": "company-1", "name": "Test", "isAdmin": true},
				},
				"paging": map[string]interface{}{
					"count": 1.0, "limit": 50.0, "offset": 0.0, "next": false,
				},
			})
		case "/api-v2/auth/keys":
			createKeyCalled = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]string{"key": "test-api-key-123"})
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	key, err := Login(context.Background(), srv.URL, "user@example.com", "secret")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if key != "test-api-key-123" {
		t.Errorf("key = %q, want test-api-key-123", key)
	}
	if !companiesCalled {
		t.Error("GetCompanies was not called")
	}
	if !createKeyCalled {
		t.Error("AuthKeyControllerCreate was not called")
	}
}

func TestLogin_GetCompanies_401_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api-v2/auth/companies" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := Login(context.Background(), srv.URL, "bad@example.com", "wrong")
	if err == nil {
		t.Fatal("expected error on 401")
	}
	if !strings.Contains(err.Error(), "get companies") {
		t.Errorf("error should mention get companies: %v", err)
	}
}

func TestLogin_NoCompanies_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api-v2/auth/companies" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"content": []interface{}{},
				"paging":  map[string]interface{}{"count": 0.0, "limit": 50.0, "offset": 0.0, "next": false},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := Login(context.Background(), srv.URL, "user@example.com", "secret")
	if err == nil {
		t.Fatal("expected error when no companies")
	}
	if !strings.Contains(err.Error(), "no companies") {
		t.Errorf("error should mention no companies: %v", err)
	}
}
