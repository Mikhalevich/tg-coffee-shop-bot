package app

import (
	"github.com/Mikhalevich/tgbot"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/bot/internal/app/tghandler"
)

func makeRoutes(tbot *tgbot.TGBot, handler *tghandler.TGHandler) {
	tbot.AddTextCommand("/start", handler.Start)

	tbot.AddMenuCommand("/order", "order food", handler.Order)
	tbot.AddMenuCommand("/order_info", "information about active order", handler.GetActiveOrder)
	tbot.AddMenuCommand("/queue", "current orders queue size", handler.OrderQueueSize)
	tbot.AddTextCommand("/history_v1", handler.OrderHistory) // obsolete
	tbot.AddMenuCommand("/history", "view history orders", handler.OrderHistoryV2)

	tbot.AddDefaultHandler(handler.DefaultHandler)
	tbot.AddDefaultCallbackQueryHandler(handler.DefaultCallbackQuery)
}
