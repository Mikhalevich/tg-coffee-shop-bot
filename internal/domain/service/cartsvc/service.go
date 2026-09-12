package cartsvc

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/cartorder"
)

var (
	_ cartorder.CartService = (*Service)(nil)
)

type Repository interface {
	StartNewCart(ctx context.Context, chatID msginfo.ChatID) (cart.ID, error)
	Clear(ctx context.Context, chatID msginfo.ChatID, cartID cart.ID) error
	AddProduct(ctx context.Context, id cart.ID, p cart.CartProduct) error
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
