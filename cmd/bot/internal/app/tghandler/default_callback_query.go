package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tgbot"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (t *TGHandler) DefaultCallbackQuery(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	btn, err := t.buttonProvider.GetButton(ctx, button.IDFromString(msg.Data))
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			sender.SendMessage(ctx, msg.ChatID, "Action already executed or expired")

			return nil
		}

		return fmt.Errorf("get button: %w", err)
	}

	if btn.ChatID.Int64() != msg.ChatID {
		return fmt.Errorf("chat not match button: %d msg: %d", btn.ChatID.Int64(), msg.ChatID)
	}

	info := msginfo.Info{
		ChatID:    msginfo.ChatIDFromInt64(msg.ChatID),
		MessageID: msginfo.MessageIDFromInt(msg.MessageID),
	}

	handler, ok := t.cbHandlers[btn.Operation]
	if !ok {
		return fmt.Errorf("operation %s is not implented", btn.Operation)
	}

	if err := handler(ctx, info, *btn); err != nil {
		return fmt.Errorf("cb handler operation %s failure: %w", btn.Operation, err)
	}

	return nil
}

func (t *TGHandler) cancelOrder(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CancelOrderPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.actionProcessor.Cancel(
		ctx,
		info.ChatID,
		info.MessageID,
		payload.OrderID,
		payload.IsTextMessage,
	); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}

func (t *TGHandler) confirmCart(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CartConfirmPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.cartProcessor.Confirm(ctx, info, payload.CartID, payload.CurrencyID); err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	return nil
}

func (t *TGHandler) cancelCart(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CartCancelPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.cartProcessor.Cancel(ctx, info, payload.CartID); err != nil {
		return fmt.Errorf("cart cancel: %w", err)
	}

	return nil
}

func (t *TGHandler) viewCategoryProducts(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CartViewCategoryProductsPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.cartProcessor.ViewCategoryProducts(
		ctx,
		info,
		payload.CartID,
		payload.CategoryID,
		payload.CurrencyID,
	); err != nil {
		return fmt.Errorf("view category products: %w", err)
	}

	return nil
}

func (t *TGHandler) viewCategories(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CartViewCategoriesPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.cartProcessor.ViewCategories(ctx, info, payload.CartID, payload.CurrencyID); err != nil {
		return fmt.Errorf("cart view categories: %w", err)
	}

	return nil
}

func (t *TGHandler) addProduct(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.CartAddProductPayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.cartProcessor.AddProduct(
		ctx,
		info,
		payload.CartID,
		payload.CategoryID,
		payload.ProductID,
		payload.CurrencyID,
	); err != nil {
		return fmt.Errorf("cart add product: %w", err)
	}

	return nil
}

func (t *TGHandler) historyPrevious(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.OrderHistoryByID](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.historyProcessor.Previous(
		ctx,
		info,
		payload.OrderID,
	); err != nil {
		return fmt.Errorf("history previous: %w", err)
	}

	return nil
}

func (t *TGHandler) historyNext(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.OrderHistoryByID](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.historyProcessor.Next(
		ctx,
		info,
		payload.OrderID,
	); err != nil {
		return fmt.Errorf("history next: %w", err)
	}

	return nil
}

func (t *TGHandler) historyFirst(ctx context.Context, info msginfo.Info, btn button.Button) error {
	if err := t.historyProcessor.First(
		ctx,
		info,
	); err != nil {
		return fmt.Errorf("history first: %w", err)
	}

	return nil
}

func (t *TGHandler) historyLast(ctx context.Context, info msginfo.Info, btn button.Button) error {
	if err := t.historyProcessor.Last(
		ctx,
		info,
	); err != nil {
		return fmt.Errorf("history last: %w", err)
	}

	return nil
}

func (t *TGHandler) historyFirstV2(ctx context.Context, info msginfo.Info, btn button.Button) error {
	if err := t.historyProcessorV2.First(
		ctx,
		info,
	); err != nil {
		return fmt.Errorf("history first: %w", err)
	}

	return nil
}

func (t *TGHandler) historyLastV2(ctx context.Context, info msginfo.Info, btn button.Button) error {
	if err := t.historyProcessorV2.Last(
		ctx,
		info,
	); err != nil {
		return fmt.Errorf("history last: %w", err)
	}

	return nil
}

func (t *TGHandler) historyPageV2(ctx context.Context, info msginfo.Info, btn button.Button) error {
	payload, err := button.GetPayload[button.OrderHistoryByPagePayload](btn)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := t.historyProcessorV2.Page(
		ctx,
		info,
		payload.Page,
	); err != nil {
		return fmt.Errorf("history page: %w", err)
	}

	return nil
}
