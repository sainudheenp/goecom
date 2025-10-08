package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sainudheenp/goecom/internal/handler"
	"github.com/sainudheenp/goecom/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx interface{}, req service.RegisterRequest) (*service.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx interface{}, req service.LoginRequest) (*service.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.LoginResponse), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful registration", func(t *testing.T) {
		mockService := new(MockAuthService)
		h := handler.NewAuthHandler(mockService)

		req := service.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			FullName: "Test User",
		}

		mockService.On("Register", mock.Anything, req).Return(&service.RegisterResponse{
			Email:    req.Email,
			FullName: req.FullName,
		}, nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := new(MockAuthService)
		h := handler.NewAuthHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
