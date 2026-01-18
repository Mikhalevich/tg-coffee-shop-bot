package order

import (
	"time"
)

type UpdateOrderData struct {
	Status              Status
	StatusOperationTime time.Time
	VerificationCode    string
	DailyPosition       int
}
