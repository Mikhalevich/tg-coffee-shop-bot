package outboxmsg

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type Invoice struct {
	order.Invoice

	ID int
}
