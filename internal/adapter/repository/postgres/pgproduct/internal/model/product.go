package model

import (
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type Product struct {
	ID         int       `db:"id"`
	Title      string    `db:"title"`
	CurrencyID int       `db:"currency_id"`
	Price      int       `db:"price"`
	IsEnabled  bool      `db:"is_enabled"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (p Product) ToDom() product.Product {
	return product.Product{
		ID:         product.ProductIDFromInt(p.ID),
		Title:      p.Title,
		CurrencyID: currency.IDFromInt(p.CurrencyID),
		Price:      p.Price,
		IsEnabled:  p.IsEnabled,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

func ToDomProducts(dbProducts []Product) []product.Product {
	portProducts := make([]product.Product, 0, len(dbProducts))

	for _, p := range dbProducts {
		portProducts = append(portProducts, p.ToDom())
	}

	return portProducts
}

type Category struct {
	ID        int    `db:"id"`
	Title     string `db:"title"`
	IsEnabled bool   `db:"is_enabled"`
}

func (c Category) ToDom() product.Category {
	return product.Category{
		ID:        product.CategoryIDFromInt(c.ID),
		Title:     c.Title,
		IsEnabled: c.IsEnabled,
	}
}

type ProductCategory struct {
	ProductID     int       `db:"product_id"`
	CategoryID    int       `db:"category_id"`
	ProductTitle  string    `db:"product_title"`
	CategoryTitle string    `db:"category_title"`
	Price         int       `db:"price"`
	IsEnabled     bool      `db:"is_enabled"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (p ProductCategory) ToDom() product.Product {
	return product.Product{
		ID:        product.ProductIDFromInt(p.ProductID),
		Title:     p.ProductTitle,
		Price:     p.Price,
		IsEnabled: p.IsEnabled,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToPortCategories(dbCategories []Category) []product.Category {
	portCategories := make([]product.Category, 0, len(dbCategories))
	for _, v := range dbCategories {
		portCategories = append(portCategories, v.ToDom())
	}

	return portCategories
}

func ToPortCategoryProducts(dbProducts []ProductCategory) []product.CategoryProducts {
	categoryMap := make(map[int]product.CategoryProducts)
	for _, dbProduct := range dbProducts {
		if portProducts, ok := categoryMap[dbProduct.CategoryID]; ok {
			portProducts.Products = append(portProducts.Products, dbProduct.ToDom())
			categoryMap[dbProduct.CategoryID] = portProducts

			continue
		}

		categoryMap[dbProduct.CategoryID] = product.CategoryProducts{
			ID:       product.CategoryIDFromInt(dbProduct.CategoryID),
			Title:    dbProduct.CategoryTitle,
			Products: []product.Product{dbProduct.ToDom()},
		}
	}

	category := make([]product.CategoryProducts, 0, len(categoryMap))
	for _, v := range categoryMap {
		category = append(category, v)
	}

	return category
}
