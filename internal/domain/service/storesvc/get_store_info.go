package storesvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

func (s *Service) GetStoreInfo(ctx context.Context) (store.StoreInfo, error) {
	stor, err := s.repo.GetStoreByID(ctx, s.storeID)
	if err != nil {
		return store.StoreInfo{}, fmt.Errorf("get store by id %d: %w", s.storeID.IntID.Int(), err)
	}

	nextWorkingTime, isActive := stor.Schedule.NextWorkingTime(s.timePrivider.Now())

	return store.StoreInfo{
		ID:                stor.ID,
		Description:       stor.Description,
		DefaultCurrencyID: stor.DefaultCurrencyID,
		IsActive:          isActive,
		NextWorkingTime:   nextWorkingTime,
	}, nil
}
