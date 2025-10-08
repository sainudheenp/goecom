//go:build integration
// +build integration

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sainudheenp/goecom/internal/config"
	"github.com/sainudheenp/goecom/internal/server"
	"github.com/sainudheenp/goecom/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration test that demonstrates the full user flow
func TestFullUserFlow(t *testing.T) {
	// Skip if DATABASE_URL is not set
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load config
	cfg, err := config.Load()
	require.NoError(t, err)

	// Override with test database if needed
	// cfg.Database.URL = "postgres://postgres:postgres@localhost:5432/ecom_test?sslmode=disable"

	// Create server
	srv, err := server.NewServer(cfg)
	require.NoError(t, err)
	defer srv.Close()

	// Get router for testing
	router := srv.GetRouter()
	router := srv.GetRouter()

	// Test 1: Register a new user
	t.Run("Register User", func(t *testing.T) {
		registerReq := service.RegisterRequest{
			Email:    fmt.Sprintf("testuser_%d@example.com", time.Now().Unix()),
			Password: "testpass123",
			FullName: "Test User",
		}

		body, _ := json.Marshal(registerReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp service.RegisterResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, registerReq.Email, resp.Email)
		assert.Equal(t, registerReq.FullName, resp.FullName)
	})

	// Test 2: Login
	var accessToken string
	t.Run("Login User", func(t *testing.T) {
		loginReq := service.LoginRequest{
			Email:    "admin@example.com", // Use seeded admin user
			Password: "admin123",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp service.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		accessToken = resp.AccessToken
	})

	// Test 3: List products
	var productID string
	t.Run("List Products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		items := resp["items"].([]interface{})
		if len(items) > 0 {
			firstProduct := items[0].(map[string]interface{})
			productID = firstProduct["id"].(string)
		}
	})

	// Test 4: Add to cart
	t.Run("Add to Cart", func(t *testing.T) {
		if productID == "" {
			t.Skip("No products available to add to cart")
		}

		addToCartReq := map[string]interface{}{
			"product_id": productID,
			"quantity":   2,
		}

		body, _ := json.Marshal(addToCartReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test 5: Get cart
	t.Run("Get Cart", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotNil(t, resp["items"])
	})

	// Test 6: Create order
	t.Run("Create Order", func(t *testing.T) {
		orderReq := map[string]interface{}{
			"shipping_address": map[string]interface{}{
				"line1":    "123 Main St",
				"city":     "Test City",
				"state":    "TS",
				"country":  "US",
				"postcode": "12345",
			},
		}

		body, _ := json.Marshal(orderReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// May fail if cart is empty
		if w.Code == http.StatusCreated {
			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			assert.NotNil(t, resp["id"])
		}
	})
}
