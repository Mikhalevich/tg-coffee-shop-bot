package productsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (s *Service) GetProductsByIDs(
	ctx context.Context,
	ids []product.ProductID,
	currencyID currency.ID,
) (map[product.ProductID]product.Product, error) {
	products, err := s.repo.GetProductsByIDs(
		ctx,
		ids,
		currencyID,
	)

	if err != nil {
		return nil, fmt.Errorf("get products from repo: %w", err)
	}

	return products, nil
}
