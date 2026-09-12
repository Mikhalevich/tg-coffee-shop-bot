package currencysvc

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/cartorder"
)

var (
	_ cartorder.CurrencyService = (*Service)(nil)
)

type Repository interface {
	GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error)
}

type Service struct {
	repo Repository
}

func New(
	repo Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}
