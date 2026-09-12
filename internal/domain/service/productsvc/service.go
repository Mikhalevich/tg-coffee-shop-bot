package productsvc

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/cart/create"
)

var (
	_ create.ProductService = (*Service)(nil)
)

type Repository interface {
	GetCategories(ctx context.Context) ([]product.Category, error)
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
