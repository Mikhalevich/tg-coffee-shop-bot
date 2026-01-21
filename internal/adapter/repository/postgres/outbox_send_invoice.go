package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (p *Postgres) SendInvoice(
	ctx context.Context,
	chatID msginfo.ChatID,
	title string,
	ord *order.Order,
	productsInfo map[product.ProductID]product.Product,
	curr *currency.Currency,
) error {
	var (
		query = `
			INSERT INTO outbox_order_invoice(
				chat_id,
				msg_text,
				order_id
			) VALUES (
				:chat_id,
				:msg_text,
				:order_id
			)
		`

		invoice = model.OutboxInvoice{
			ChatID:  chatID.Int64(),
			Text:    title,
			OrderID: ord.ID.Int(),
		}
	)

	res, err := sqlx.NamedExecContext(ctx, p.Transactor().ExtContext(ctx), query, invoice)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
