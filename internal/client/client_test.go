package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	baseURL := "https://api.example.com"

	client := NewClient(baseURL, apiKey)

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseURL %s, got %s", baseURL, client.BaseURL)
	}

	if client.APIKey != apiKey {
		t.Errorf("Expected APIKey %s, got %s", apiKey, client.APIKey)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized")
	}
}

func TestDoRequest_Authentication(t *testing.T) {
	apiKey := "test-api-key"

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("X-Api-Key")

		if authHeader != apiKey {
			t.Errorf("Expected X-Api-Key header %s, got %s", apiKey, authHeader)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, apiKey)

	// Make a test request
	err := client.doRequest("GET", "", nil, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestDoRequest_ErrorHandling(t *testing.T) {
	// Create test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message": "Resource not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	err := client.doRequest("GET", "", nil, nil)
	if err == nil {
		t.Error("Expected error for 404 response")
	}
}
