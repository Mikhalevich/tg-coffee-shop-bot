package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (p *Postgres) OutboxSendMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	textType messageprocessor.MessageTextType,
	rows ...button.ButtonRow,
) error {
	if err := p.insertOutboxMessage(ctx, model.OutboxMessage{
		ChatID:    chatID.Int64(),
		Text:      text,
		Type:      model.ToDBMessageType(textType),
		CreatedAt: time.Now(),
	}); err != nil {
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
