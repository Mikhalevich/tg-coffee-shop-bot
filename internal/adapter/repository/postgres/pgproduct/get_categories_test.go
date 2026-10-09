package pgproduct_test

import (
	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (s *ProductSuit) TestGetCategories() {
	s.Run("enabled categories with enabled products ordered by title", func() {
		var (
			ctx        = s.T().Context()
			tea        = s.insertCategory("Tea", true)
			coffee     = s.insertCategory("Coffee", true)
			desserts   = s.insertCategory("Desserts", true)
			latte      = s.insertProduct("Latte", true)
			espresso   = s.insertProduct("Espresso", true)
			greenTea   = s.insertProduct("Green tea", true)
			cheesecake = s.insertProduct("Cheesecake", true)
		)

		s.linkProductCategory(latte, coffee)
		s.linkProductCategory(espresso, coffee)
		s.linkProductCategory(greenTea, tea)
		s.linkProductCategory(cheesecake, desserts)

		actual, err := s.pgProduct.GetCategories(ctx)

		s.Require().NoError(err)
		s.Require().Equal([]product.Category{
			{ID: product.CategoryIDFromInt(coffee), Title: "Coffee", IsEnabled: true},
			{ID: product.CategoryIDFromInt(desserts), Title: "Desserts", IsEnabled: true},
			{ID: product.CategoryIDFromInt(tea), Title: "Tea", IsEnabled: true},
		}, actual)
	})

	s.Run("skip disabled category", func() {
		var (
			ctx      = s.T().Context()
			coffee   = s.insertCategory("Coffee", true)
			seasonal = s.insertCategory("Seasonal", false)
			latte    = s.insertProduct("Latte", true)
			pumpkin  = s.insertProduct("Pumpkin latte", true)
		)

		s.linkProductCategory(latte, coffee)
		s.linkProductCategory(pumpkin, seasonal)

		actual, err := s.pgProduct.GetCategories(ctx)

		s.Require().NoError(err)
		s.Require().Equal([]product.Category{
			{ID: product.CategoryIDFromInt(coffee), Title: "Coffee", IsEnabled: true},
		}, actual)
	})

	s.Run("skip category with only disabled products", func() {
		var (
			ctx      = s.T().Context()
			coffee   = s.insertCategory("Coffee", true)
			desserts = s.insertCategory("Desserts", true)
			latte    = s.insertProduct("Latte", true)
			brownie  = s.insertProduct("Brownie", false)
		)

		s.linkProductCategory(latte, coffee)
		s.linkProductCategory(brownie, desserts)

		actual, err := s.pgProduct.GetCategories(ctx)

		s.Require().NoError(err)
		s.Require().Equal([]product.Category{
			{ID: product.CategoryIDFromInt(coffee), Title: "Coffee", IsEnabled: true},
		}, actual)
	})

	s.Run("skip category without products", func() {
		var (
			ctx    = s.T().Context()
			coffee = s.insertCategory("Coffee", true)
			latte  = s.insertProduct("Latte", true)
		)

		s.insertCategory("Empty", true)
		s.linkProductCategory(latte, coffee)

		actual, err := s.pgProduct.GetCategories(ctx)

		s.Require().NoError(err)
		s.Require().Equal([]product.Category{
			{ID: product.CategoryIDFromInt(coffee), Title: "Coffee", IsEnabled: true},
		}, actual)
	})

	s.Run("product in multiple categories", func() {
		var (
			ctx       = s.T().Context()
			coffee    = s.insertCategory("Coffee", true)
			breakfast = s.insertCategory("Breakfast", true)
			latte     = s.insertProduct("Latte", true)
		)

		s.linkProductCategory(latte, coffee)
		s.linkProductCategory(latte, breakfast)

		actual, err := s.pgProduct.GetCategories(ctx)

		s.Require().NoError(err)
		s.Require().Equal([]product.Category{
			{ID: product.CategoryIDFromInt(breakfast), Title: "Breakfast", IsEnabled: true},
			{ID: product.CategoryIDFromInt(coffee), Title: "Coffee", IsEnabled: true},
		}, actual)
	})

	s.Run("empty", func() {
		ctx := s.T().Context()

		s.insertCategory("Coffee", false)
		s.insertProduct("Latte", true)

		actual, err := s.pgProduct.GetCategories(ctx)

		s.Require().NoError(err)
		s.Require().Empty(actual)
	})
}

func (s *ProductSuit) insertCategory(title string, isEnabled bool) int {
	var (
		ctx        = s.T().Context()
		categoryID int
	)

	err := sqlx.GetContext(ctx, s.transactor.ExtContext(ctx), &categoryID, `
		INSERT INTO category (
			title,
			is_enabled
		) VALUES (
			$1,
			$2
		) RETURNING id`,
		title,
		isEnabled,
	)
	s.Require().NoError(err)

	return categoryID
}

func (s *ProductSuit) insertProduct(title string, isEnabled bool) int {
	var (
		ctx       = s.T().Context()
		productID int
	)

	err := sqlx.GetContext(ctx, s.transactor.ExtContext(ctx), &productID, `
		INSERT INTO product (
			title,
			is_enabled,
			created_at,
			updated_at
		) VALUES (
			$1,
			$2,
			NOW(),
			NOW()
		) RETURNING id`,
		title,
		isEnabled,
	)
	s.Require().NoError(err)

	return productID
}

func (s *ProductSuit) linkProductCategory(productID int, categoryID int) {
	ctx := s.T().Context()

	_, err := s.transactor.ExtContext(ctx).ExecContext(ctx, `
		INSERT INTO product_category (
			product_id,
			category_id
		) VALUES (
			$1,
			$2
		)`,
		productID,
		categoryID,
	)
	s.Require().NoError(err)
}
