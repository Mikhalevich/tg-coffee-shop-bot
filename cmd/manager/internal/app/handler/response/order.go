package response

import (
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type Order struct {
	ID               string       `json:"id" example:"123" doc:"Order id"`
	Status           string       `json:"status" example:"created" doc:"Order status"`
	VerificationCode string       `json:"verification_code" example:"012" doc:"Order verification code"`
	Timeline         []StatusTime `json:"timeline" doc:"Timeline status changes"`
}

type StatusTime struct {
	Status string    `json:"status" example:"created" doc:"Order status"`
	Time   time.Time `json:"time" doc:"Status time"`
}

func ToOrder(ord order.Order) *Order {
	timeline := make([]StatusTime, 0, len(ord.Timeline))
	for _, tl := range ord.Timeline {
		timeline = append(timeline, StatusTime{
			Status: tl.Status.String(),
			Time:   tl.Time,
		})
	}

	return &Order{
		ID:               ord.ID.String(),
		Status:           ord.Status.String(),
		VerificationCode: ord.VerificationCode,
		Timeline:         timeline,
	}
}
