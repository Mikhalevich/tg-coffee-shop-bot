package pgoutbox

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgoutbox/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOutbox) SendMessage(
	ctx context.Context,
	msg msginfo.Message,
) error {
	dbMsg, err := model.ToDBOutboxMessage(msg)
	if err != nil {
		return fmt.Errorf("make db outbox message: %w", err)
	}

	var (
		query = `
			INSERT INTO outbox_messages(
				chat_id,
				reply_msg_id,
				msg_text,
				msg_type,
				payload,
				buttons
			) VALUES (
				:chat_id,
				:reply_msg_id,
				:msg_text,
				:msg_type,
				:payload,
				:buttons
			)
		`
	)

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		dbMsg,
	)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return perror.NoRowsUpdated()
	}

	return nil
}
