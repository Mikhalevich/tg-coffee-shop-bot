package outboxmsg

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type Invoice struct {
	ID      int
	ChatID  msginfo.ChatID
	Text    string
	OrderID order.ID
}
