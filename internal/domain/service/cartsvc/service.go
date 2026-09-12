package cartsvc

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/cart/create"
)

var (
	_ create.CartService = (*Service)(nil)
)

type Repository interface {
	StartNewCart(ctx context.Context, chatID msginfo.ChatID) (cart.ID, error)
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
