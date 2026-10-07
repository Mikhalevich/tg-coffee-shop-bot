package cartsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
)

func (s *Service) AddProduct(
	ctx context.Context,
	id cart.ID,
	p cart.CartProduct,
) error {
	if err := s.repo.AddProduct(ctx, id, p); err != nil {
		return fmt.Errorf("add product: %w", err)
	}

	return nil
}
