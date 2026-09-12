package pgcurrency_test

import (
	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
)

func (s *CurrencySuit) TestGetCurrencyByID() {
	s.Run("success", func() {
		var (
			ctx      = s.T().Context()
			expected = &currency.Currency{
				ID:         currency.IDFromInt(1),
				Code:       "USD",
				Exp:        2,
				DecimalSep: ".",
				MinAmount:  100,
				MaxAmount:  100000,
				IsEnabled:  true,
			}
		)

		_, err := sqlx.NamedExecContext(ctx, s.transactor.ExtContext(ctx), `
			INSERT INTO currency (
				code,
				exp,
				decimal_sep,
				min_amount,
				max_amount,
				is_enabled
			) VALUES (
				:code,
				:exp,
				:decimal_sep,
				:min_amount,
				:max_amount,
				:is_enabled
			)`, map[string]any{
			"code":        expected.Code,
			"exp":         expected.Exp,
			"decimal_sep": expected.DecimalSep,
			"min_amount":  expected.MinAmount,
			"max_amount":  expected.MaxAmount,
			"is_enabled":  expected.IsEnabled,
		})
		s.Require().NoError(err)

		actual, err := s.pgCurrency.GetCurrencyByID(ctx, expected.ID)

		s.Require().NoError(err)
		s.Require().Equal(expected, actual)
	})

	s.Run("not found", func() {
		var (
			ctx = s.T().Context()
			id  = currency.IDFromInt(999)
		)

		actual, err := s.pgCurrency.GetCurrencyByID(ctx, id)

		s.Require().Error(err)
		s.Require().Nil(actual)
		s.Require().EqualError(err, "currency not found")
	})
}
