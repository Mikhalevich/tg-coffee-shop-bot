package pgoutbox

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgoutbox/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOutbox) SendInvoice(
	ctx context.Context,
	invoice order.Invoice,
) error {
	var (
		query = `
			INSERT INTO outbox_order_invoice(
				chat_id,
				title,
				description,
				currency_code,
				order_id,
				labels,
				buttons
			) VALUES (
				:chat_id,
				:title,
				:description,
				:currency_code,
				:order_id,
				:labels,
				:buttons
			)
		`
	)

	dbInvoice, err := model.ToDBInvoice(invoice)
	if err != nil {
		return fmt.Errorf("convert to db invoice: %w", err)
	}

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		dbInvoice,
	)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return perror.NoRowsUpdated()
	}

	return nil
}
