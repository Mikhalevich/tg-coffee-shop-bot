package pgproduct_test

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (s *ProductSuit) TestGetProductsByCategoryID() {
	s.Run("enabled products of category ordered by title", func() {
		var (
			ctx      = s.T().Context()
			usd      = s.insertCurrency("USD")
			coffee   = s.insertCategory("Coffee", true)
			latte    = s.insertProduct("Latte", true)
			espresso = s.insertProduct("Espresso", true)
		)

		s.linkProductCategory(latte, coffee)
		s.linkProductCategory(espresso, coffee)
		s.insertPrice(latte, usd, 350)
		s.insertPrice(espresso, usd, 250)

		actual, err := s.pgProduct.GetProductsByCategoryID(
			ctx,
			product.CategoryIDFromInt(coffee),
			currency.IDFromInt(usd),
		)

		s.Require().NoError(err)
		s.Require().Len(actual, 2)
		s.Require().Equal([]product.Product{
			makeProduct(espresso, "Espresso", usd, 250, true),
			makeProduct(latte, "Latte", usd, 350, true),
		}, []product.Product{
			normalizeTime(actual[0]),
			normalizeTime(actual[1]),
		})
	})

	s.Run("skip disabled products", func() {
		var (
			ctx      = s.T().Context()
			usd      = s.insertCurrency("USD")
			coffee   = s.insertCategory("Coffee", true)
			latte    = s.insertProduct("Latte", true)
			espresso = s.insertProduct("Espresso", false)
		)

		s.linkProductCategory(latte, coffee)
		s.linkProductCategory(espresso, coffee)
		s.insertPrice(latte, usd, 350)
		s.insertPrice(espresso, usd, 250)

		actual, err := s.pgProduct.GetProductsByCategoryID(
			ctx,
			product.CategoryIDFromInt(coffee),
			currency.IDFromInt(usd),
		)

		s.Require().NoError(err)
		s.Require().Len(actual, 1)
		s.Require().Equal(makeProduct(latte, "Latte", usd, 350, true), normalizeTime(actual[0]))
	})

	s.Run("skip products of other categories", func() {
		var (
			ctx      = s.T().Context()
			usd      = s.insertCurrency("USD")
			coffee   = s.insertCategory("Coffee", true)
			desserts = s.insertCategory("Desserts", true)
			latte    = s.insertProduct("Latte", true)
			brownie  = s.insertProduct("Brownie", true)
		)

		s.linkProductCategory(latte, coffee)
		s.linkProductCategory(brownie, desserts)
		s.insertPrice(latte, usd, 350)
		s.insertPrice(brownie, usd, 400)

		actual, err := s.pgProduct.GetProductsByCategoryID(
			ctx,
			product.CategoryIDFromInt(coffee),
			currency.IDFromInt(usd),
		)

		s.Require().NoError(err)
		s.Require().Len(actual, 1)
		s.Require().Equal(makeProduct(latte, "Latte", usd, 350, true), normalizeTime(actual[0]))
	})

	s.Run("price in requested currency", func() {
		var (
			ctx      = s.T().Context()
			usd      = s.insertCurrency("USD")
			eur      = s.insertCurrency("EUR")
			coffee   = s.insertCategory("Coffee", true)
			latte    = s.insertProduct("Latte", true)
			espresso = s.insertProduct("Espresso", true)
		)

		s.linkProductCategory(latte, coffee)
		s.linkProductCategory(espresso, coffee)
		s.insertPrice(latte, usd, 350)
		s.insertPrice(latte, eur, 320)
		s.insertPrice(espresso, usd, 250)

		actual, err := s.pgProduct.GetProductsByCategoryID(
			ctx,
			product.CategoryIDFromInt(coffee),
			currency.IDFromInt(eur),
		)

		s.Require().NoError(err)
		s.Require().Len(actual, 1)
		s.Require().Equal(makeProduct(latte, "Latte", eur, 320, true), normalizeTime(actual[0]))
	})

	s.Run("empty", func() {
		var (
			ctx    = s.T().Context()
			usd    = s.insertCurrency("USD")
			coffee = s.insertCategory("Coffee", true)
		)

		actual, err := s.pgProduct.GetProductsByCategoryID(
			ctx,
			product.CategoryIDFromInt(coffee),
			currency.IDFromInt(usd),
		)

		s.Require().NoError(err)
		s.Require().Empty(actual)
	})
}
