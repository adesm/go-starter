package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	Update(ctx context.Context, id uint, req *UpdateUserRequest) (*User, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, pageSize int) ([]*User, error)
}
