package create

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type StoreService interface {
	GetStoreInfo(ctx context.Context) (store.StoreInfo, error)
}

type Create struct {
	storeService StoreService
}

func New(
	storeService StoreService,
) *Create {
	return &Create{
		storeService: storeService,
	}
}

func (c *Create) Create(ctx context.Context) error {
	return nil
}
