package model

import (
	"fmt"
	"sort"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type OrderTimeline struct {
	ID        int       `db:"order_id"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}

func toDomTimeline(dbTimeline []OrderTimeline) ([]order.StatusTime, error) {
	sort.Slice(dbTimeline, func(i, j int) bool {
		return dbTimeline[i].UpdatedAt.Sub(dbTimeline[j].UpdatedAt) < 0
	})

	portTimeline := make([]order.StatusTime, 0, len(dbTimeline))

	for _, orderTime := range dbTimeline {
		status, err := order.StatusFromString(orderTime.Status)
		if err != nil {
			return nil, fmt.Errorf("timeline status from string: %w", err)
		}

		portTimeline = append(portTimeline, order.StatusTime{
			Status: status,
			Time:   orderTime.UpdatedAt,
		})
	}

	return portTimeline, nil
}
