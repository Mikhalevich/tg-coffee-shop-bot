package currencysvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
)

func (s *Service) GetCurrencyByID(
	ctx context.Context,
	id currency.ID,
) (currency.Currency, error) {
	curr, err := s.repo.GetCurrencyByID(ctx, id)
	if err != nil {
		return currency.Currency{}, fmt.Errorf("get currency by id: %w", err)
	}

	return curr, nil
}
