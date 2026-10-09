package pgproduct_test

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (s *ProductSuit) TestGetProductsByIDs() {
	s.Run("success", func() {
		var (
			ctx      = s.T().Context()
			usd      = s.insertCurrency("USD")
			latte    = s.insertProduct("Latte", true)
			espresso = s.insertProduct("Espresso", true)
		)

		s.insertPrice(latte, usd, 350)
		s.insertPrice(espresso, usd, 250)

		actual, err := s.pgProduct.GetProductsByIDs(
			ctx,
			[]product.ProductID{product.ProductIDFromInt(latte), product.ProductIDFromInt(espresso)},
			currency.IDFromInt(usd),
		)

		s.Require().NoError(err)
		s.Require().Equal(map[product.ProductID]product.Product{
			product.ProductIDFromInt(latte):    makeProduct(latte, "Latte", usd, 350, true),
			product.ProductIDFromInt(espresso): makeProduct(espresso, "Espresso", usd, 250, true),
		}, normalizeTimeMap(actual))
	})

	s.Run("only requested ids", func() {
		var (
			ctx      = s.T().Context()
			usd      = s.insertCurrency("USD")
			latte    = s.insertProduct("Latte", true)
			espresso = s.insertProduct("Espresso", true)
		)

		s.insertPrice(latte, usd, 350)
		s.insertPrice(espresso, usd, 250)

		actual, err := s.pgProduct.GetProductsByIDs(
			ctx,
			[]product.ProductID{product.ProductIDFromInt(latte)},
			currency.IDFromInt(usd),
		)

		s.Require().NoError(err)
		s.Require().Equal(map[product.ProductID]product.Product{
			product.ProductIDFromInt(latte): makeProduct(latte, "Latte", usd, 350, true),
		}, normalizeTimeMap(actual))
	})

	s.Run("disabled products are returned", func() {
		var (
			ctx   = s.T().Context()
			usd   = s.insertCurrency("USD")
			latte = s.insertProduct("Latte", false)
		)

		s.insertPrice(latte, usd, 350)

		actual, err := s.pgProduct.GetProductsByIDs(
			ctx,
			[]product.ProductID{product.ProductIDFromInt(latte)},
			currency.IDFromInt(usd),
		)

		s.Require().NoError(err)
		s.Require().Equal(map[product.ProductID]product.Product{
			product.ProductIDFromInt(latte): makeProduct(latte, "Latte", usd, 350, false),
		}, normalizeTimeMap(actual))
	})

	s.Run("price in requested currency", func() {
		var (
			ctx      = s.T().Context()
			usd      = s.insertCurrency("USD")
			eur      = s.insertCurrency("EUR")
			latte    = s.insertProduct("Latte", true)
			espresso = s.insertProduct("Espresso", true)
		)

		s.insertPrice(latte, usd, 350)
		s.insertPrice(latte, eur, 320)
		s.insertPrice(espresso, usd, 250)

		actual, err := s.pgProduct.GetProductsByIDs(
			ctx,
			[]product.ProductID{product.ProductIDFromInt(latte), product.ProductIDFromInt(espresso)},
			currency.IDFromInt(eur),
		)

		s.Require().NoError(err)
		s.Require().Equal(map[product.ProductID]product.Product{
			product.ProductIDFromInt(latte): makeProduct(latte, "Latte", eur, 320, true),
		}, normalizeTimeMap(actual))
	})

	s.Run("not found", func() {
		var (
			ctx = s.T().Context()
			usd = s.insertCurrency("USD")
		)

		actual, err := s.pgProduct.GetProductsByIDs(
			ctx,
			[]product.ProductID{product.ProductIDFromInt(999)},
			currency.IDFromInt(usd),
		)

		s.Require().NoError(err)
		s.Require().Nil(actual)
	})

	s.Run("empty ids", func() {
		var (
			ctx = s.T().Context()
			usd = s.insertCurrency("USD")
		)

		actual, err := s.pgProduct.GetProductsByIDs(
			ctx,
			[]product.ProductID{},
			currency.IDFromInt(usd),
		)

		s.Require().Error(err)
		s.Require().Nil(actual)
	})
}

func normalizeTimeMap(products map[product.ProductID]product.Product) map[product.ProductID]product.Product {
	normalized := make(map[product.ProductID]product.Product, len(products))

	for id, p := range products {
		normalized[id] = normalizeTime(p)
	}

	return normalized
}
