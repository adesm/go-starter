package order

import "context"

type Repository interface {
	Create(ctx context.Context, item *Order) error
	FindByID(ctx context.Context, id uint) (*Order, error)
	Update(ctx context.Context, item *Order) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*Order, error)
}

type Service interface {
	Create(ctx context.Context, req *CreateOrderRequest) (*Order, error)
	GetByID(ctx context.Context, id uint) (*Order, error)
	Update(ctx context.Context, id uint, req *UpdateOrderRequest) (*Order, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, pageSize int) ([]*Order, error)
}
