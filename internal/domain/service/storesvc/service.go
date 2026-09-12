package storesvc

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/cartorder"
)

var (
	_ cartorder.StoreService = (*Service)(nil)
)

type Repository interface {
	GetStoreByID(ctx context.Context, id store.ID) (*store.Store, error)
}

type TimeProvider interface {
	Now() time.Time
}

type Service struct {
	storeID      store.ID
	repo         Repository
	timePrivider TimeProvider
}

func New(
	storeID store.ID,
	repo Repository,
	timeProvider TimeProvider,
) *Service {
	return &Service{
		storeID:      storeID,
		repo:         repo,
		timePrivider: timeProvider,
	}
}
