package order

import (
	"context"

	"boilerplate/internal/config"
)

type service struct {
	repo Repository
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{repo: repo, cfg: cfg}
}

func (s *service) Create(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
	item := &Order{
		// Map fields from req
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*Order, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uint, req *UpdateOrderRequest) (*Order, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Update fields
	
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, page, pageSize int) ([]*Order, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, pageSize, offset)
}
