package pgbutton

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/messagesvc"
)

var (
	_ messagesvc.ButtonRepository = (*PgButton)(nil)
)

type Transactor interface {
	ExtContext(ctx context.Context) sqlx.ExtContext
}

type PgButton struct {
	transactor Transactor
}

func New(
	transactor Transactor,
) *PgButton {
	return &PgButton{
		transactor: transactor,
	}
}
