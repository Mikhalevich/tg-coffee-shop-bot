package storesvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

func (s *Service) GetStoreInfo(ctx context.Context) (store.StoreInfo, error) {
	stor, err := s.repo.GetStoreByID(ctx, s.storeID)
	if err != nil {
		return store.StoreInfo{}, fmt.Errorf("get store by id %d: %w", s.storeID.Int(), err)
	}

	now := s.timePrivider.Now()

	nextWorkingTime, isActive := stor.Schedule.NextWorkingTime(now)

	return store.StoreInfo{
		ID:                stor.ID,
		Description:       stor.Description,
		DefaultCurrencyID: stor.DefaultCurrencyID,
		IsActive:          isActive,
		CurrentTime:       now,
		NextWorkingTime:   nextWorkingTime,
	}, nil
}
