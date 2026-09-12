package cartsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (s *Service) Clear(
	ctx context.Context,
	chatID msginfo.ChatID,
	cartID cart.ID,
) error {
	if err := s.repo.Clear(ctx, chatID, cartID); err != nil {
		return fmt.Errorf("repo clear: %w", err)
	}

	return nil
}
