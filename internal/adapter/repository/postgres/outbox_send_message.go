package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
)

func (p *Postgres) SendMessage(
	ctx context.Context,
	msg messageprocessor.Message,
) error {
	if err := p.insertOutboxMessage(ctx, model.MessageToOutboxMessage(msg)); err != nil {
		return fmt.Errorf("insert outbox message: %w", err)
	}

	return nil
}

func (p *Postgres) insertOutboxMessage(ctx context.Context, msg model.OutboxMessage) error {
	var (
		query = `
			INSERT INTO outbox_messages(
				chat_id,
				reply_msg_id,
				msg_text,
				msg_type,
				payload,
				buttons,
				created_at
			) VALUES (
				:chat_id,
				:reply_msg_id,
				:msg_text,
				:msg_type,
				:payload,
				:buttons,
				:created_at
			)
		`
	)

	res, err := sqlx.NamedExecContext(ctx, p.transactor.ExtContext(ctx), query, msg)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	if _, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	return nil
}
