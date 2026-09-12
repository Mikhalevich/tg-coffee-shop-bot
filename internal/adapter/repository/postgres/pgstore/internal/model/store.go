package model

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type Store struct {
	ID                int    `db:"id"`
	Description       string `db:"description"`
	DefaultCurrencyID int    `db:"default_currency_id"`
}

type StoreSchedule struct {
	StoreID   int       `db:"store_id"`
	DayOfWeek string    `db:"day_of_week"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
}

func (s *Store) ToDom(schedule []StoreSchedule) (*store.Store, error) {
	days := make([]store.DaySchedule, 0, len(schedule))

	for _, daySchedule := range schedule {
		day, err := store.WeekdayFromString(daySchedule.DayOfWeek)
		if err != nil {
			return nil, fmt.Errorf("weekday from string: %w", err)
		}

		days = append(days, store.DaySchedule{
			Weekday:   day,
			StartTime: daySchedule.StartTime,
			EndTime:   daySchedule.EndTime,
		})
	}

	return &store.Store{
		ID:                store.IDFromInt(s.ID),
		Description:       s.Description,
		DefaultCurrencyID: currency.IDFromInt(s.DefaultCurrencyID),
		Schedule: store.Schedule{
			Days: days,
		},
	}, nil
}
