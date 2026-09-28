package pgoutbox

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/outbox/outboxsvc"
)

var (
	_ outboxsvc.Repository = (*PgOutbox)(nil)
)

type Transactor interface {
	ExtContext(ctx context.Context) sqlx.ExtContext
}

type PgOutbox struct {
	transactor Transactor
}

func New(
	transactor Transactor,
) *PgOutbox {
	return &PgOutbox{
		transactor: transactor,
	}
}
