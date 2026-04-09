package user

import (
	"context"
	"testing"
	"time"

	"boilerplate/internal/config"
	"boilerplate/internal/shared/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Register_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	cfg := &config.Config{}

	service := NewService(mockRepo, cfg)

	req := &RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock expectations
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, ErrNotFound)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *User) bool {
		return u.Email == req.Email && u.Name == req.Name
	})).Return(nil)

	user, err := service.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
	assert.True(t, utils.CheckPassword(user.Password, req.Password) == nil, "Password should be hashed properly")

	mockRepo.AssertExpectations(t)
}

func TestService_Register_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockRepository)
	cfg := &config.Config{}

	service := NewService(mockRepo, cfg)

	req := &RegisterRequest{
		Name:     "Test User",
		Email:    "existing@example.com",
		Password: "password123",
	}

	// Mock expectations
	existingUser := &User{
		ID:    1,
		Email: "existing@example.com",
	}
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	user, err := service.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrDuplicateEmail, err)

	mockRepo.AssertExpectations(t)
}

func TestService_Login_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "very-secret",
			ExpireTime: 1,
		},
	}

	service := NewService(mockRepo, cfg)

	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := utils.HashPassword("password123")
	existingUser := &User{
		ID:        1,
		Email:     "test@example.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	res, err := service.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, existingUser.Email, res.User.Email)

	mockRepo.AssertExpectations(t)
}
