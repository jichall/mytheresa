package database

import (
	"sort"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Order defines the ordering of filters and it's used when sorting filters.
type Order int

// Filter is a scope (in the GORM sense) that applies to a GORM query.
type Filter interface {
	Apply(*gorm.DB) *gorm.DB
	Ordering() Order
}

func Sort(filters []Filter) {
	sort.SliceStable(filters, func(i, j int) bool {
		return filters[i].Ordering() < filters[j].Ordering()
	})
}

type PriceFilter struct {
	Price    decimal.Decimal
	Operator string
}

func (pf *PriceFilter) Apply(db *gorm.DB) *gorm.DB {
	switch pf.Operator {
	case "gt":
		return db.Where("price > ?", pf.Price)
	case "gte":
		return db.Where("price >= ?", pf.Price)
	case "lt":
		return db.Where("price < ?", pf.Price)
	case "lte":
		return db.Where("price <= ?", pf.Price)
	case "eq":
		return db.Where("price = ?", pf.Price)
	default:
		return db
	}
}

func (pf *PriceFilter) Ordering() Order {
	return 1
}

type CategoryFilter struct {
	Categories []string
}

func (cf *CategoryFilter) Apply(db *gorm.DB) *gorm.DB {
	if len(cf.Categories) == 0 {
		return db
	}

	return db.Where("category IN (?)", cf.Categories)
}

func (cf *CategoryFilter) Ordering() Order {
	return 1
}

type PageFilter struct {
	Page int
	Size int
}

func (pf *PageFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Offset(pf.Page * pf.Size).Limit(pf.Size)
}

func (pf *PageFilter) Ordering() Order {
	return 0
}
