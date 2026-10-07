package cartsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (s *Service) Create(
	ctx context.Context,
	chatID msginfo.ChatID,
) (cart.ID, error) {
	id, err := s.repo.StartNewCart(ctx, chatID)
	if err != nil {
		return cart.IDFromString(""), fmt.Errorf("create cart: %w", err)
	}

	return id, nil
}
