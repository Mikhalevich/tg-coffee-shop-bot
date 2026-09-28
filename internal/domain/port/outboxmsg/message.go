package outboxmsg

import "github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"

type Message struct {
	msginfo.Message

	ID int
}
