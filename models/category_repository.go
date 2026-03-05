package models

import (
	"context"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (c *CategoryRepository) GetAll(ctx context.Context) ([]Category, error) {
	var categories []Category

	if err := c.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *CategoryRepository) GetPaged(ctx context.Context, page, size int) ([]Category, error) {
	var categories []Category

	total := int64(0)
	tc := &total

	if err := c.db.WithContext(ctx).Model(&Category{}).Count(tc).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * size

	err := c.db.Limit(size).Offset(offset).Order("ID desc").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *CategoryRepository) GetByCode(ctx context.Context, code string) (*Category, error) {
	category := &Category{}

	err := c.db.WithContext(ctx).Where("code = ?", code).Find(category).Error
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c *CategoryRepository) Save(ctx context.Context, category *Category) error {
	return c.db.WithContext(ctx).Save(category).Error
}
