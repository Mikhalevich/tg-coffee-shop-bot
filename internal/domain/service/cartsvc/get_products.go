package cartsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
)

func (s *Service) GetProducts(
	ctx context.Context,
	id cart.ID,
) ([]cart.CartProduct, error) {
	products, err := s.repo.GetProducts(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get products: %w", err)
	}

	return products, nil
}
