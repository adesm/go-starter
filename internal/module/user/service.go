package user

import (
	"context"

	"boilerplate/internal/config"
	"boilerplate/internal/shared/utils"
)

type service struct {
	repo Repository
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{repo: repo, cfg: cfg}
}

func (s *service) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	existing, _ := s.repo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrDuplicateEmail
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{Email: req.Email, Name: req.Name, Password: hashedPassword}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.ExpireTime)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token, User: user}, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uint, req *UpdateUserRequest) (*User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Email != user.Email {
		existing, _ := s.repo.FindByEmail(ctx, req.Email)
		if existing != nil && existing.ID != id {
			return nil, ErrDuplicateEmail
		}
	}

	user.Name = req.Name
	user.Email = req.Email

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, page, pageSize int) ([]*User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, pageSize, offset)
}
