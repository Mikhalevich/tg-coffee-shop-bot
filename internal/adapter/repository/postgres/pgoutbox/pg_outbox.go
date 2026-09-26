package pgoutbox

import (
	"context"

	"github.com/jmoiron/sqlx"
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
