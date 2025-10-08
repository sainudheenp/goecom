package service_test

import (
	"context"
	"testing"

	"github.com/sainudheenp/goecom/internal/service"
	"github.com/sainudheenp/goecom/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *store.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*store.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id interface{}) (*store.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name    string
		request service.RegisterRequest
		wantErr bool
	}{
		{
			name: "successful registration",
			request: service.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			request: service.RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
				FullName: "Test User",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			
			if tt.name == "successful registration" {
				mockRepo.On("Exists", mock.Anything, tt.request.Email).Return(false, nil)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*store.User")).Return(nil)
			}

			authService := service.NewAuthService(mockRepo, "test-secret-key-for-testing-purposes", 24, 10)
			
			_, err := authService.Register(context.Background(), tt.request)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
