package order

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, item *Order) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Order, error) {
	var item Order
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) Update(ctx context.Context, item *Order) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Order{}, id).Error
}

func (r *repository) List(ctx context.Context, limit, offset int) ([]*Order, error) {
	var items []*Order
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&items).Error
	return items, err
}
