package store

import (
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
)

type StoreInfo struct {
	ID                ID
	Description       string
	DefaultCurrencyID currency.ID
	IsActive          bool
	NextWorkingTime   time.Time
}
