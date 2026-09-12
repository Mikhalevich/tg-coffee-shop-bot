package create

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type StoreService interface {
	GetStoreInfo(ctx context.Context) (store.StoreInfo, error)
}

type NotificationService interface {
	SendStoreClosed(
		ctx context.Context,
		chatID msginfo.ChatID,
		currentTime time.Time,
		nextWorkingTime time.Time,
	) error
}

type Create struct {
	storeService        StoreService
	notificationService NotificationService
}

func New(
	storeService StoreService,
	notificationService NotificationService,
) *Create {
	return &Create{
		storeService:        storeService,
		notificationService: notificationService,
	}
}

func (c *Create) Create(
	ctx context.Context,
	info msginfo.Info,
) error {
	storeInfo, err := c.storeService.GetStoreInfo(ctx)
	if err != nil {
		return fmt.Errorf("get store info: %w", err)
	}

	if !storeInfo.IsActive {
		if err := c.notificationService.SendStoreClosed(
			ctx,
			info.ChatID,
			storeInfo.CurrentTime,
			storeInfo.NextWorkingTime,
		); err != nil {
			return fmt.Errorf("send store closed msg: %w", err)
		}
	}

	return nil
}
