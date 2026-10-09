package pgproduct_test

import (
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
