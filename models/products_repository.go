package models

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAll(ctx context.Context) ([]Product, error) {
	var products []Product

	if err := r.db.Preload("Variants").Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductsRepository) GetPaged(ctx context.Context, page, size int) ([]Product, error) {
	var products []Product

	total := int64(0)
	tc := &total

	if err := r.db.WithContext(ctx).Model(&Product{}).Count(tc).Error; err != nil {
		return nil, err
	}

	// it's indexed at 0
	offset := page * size

	err := r.db.Preload("Variants").Preload("Category").Limit(size).Offset(offset).Order("ID desc").Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductsRepository) GetByCode(ctx context.Context, code string) (*Product, error) {
	product := &Product{}

	err := r.db.WithContext(ctx).Preload("Variants").Preload("Category").Where("code = ?", code).Find(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (c *ProductsRepository) Save(ctx context.Context, product *Product) error {
	return errors.New("not implemented")
}
