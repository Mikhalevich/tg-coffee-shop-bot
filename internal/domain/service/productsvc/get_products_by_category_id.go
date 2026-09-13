package productsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (s *Service) GetProductsByCategoryID(
	ctx context.Context,
	categoryID product.CategoryID,
	currencyID currency.ID,
) ([]product.Product, error) {
	products, err := s.repo.GetProductsByCategoryID(
		ctx,
		categoryID,
		currencyID,
	)

	if err != nil {
		return nil, fmt.Errorf("get products from repo: %w", err)
	}

	return products, nil
}
