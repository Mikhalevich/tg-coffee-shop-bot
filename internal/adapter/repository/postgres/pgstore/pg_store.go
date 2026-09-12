package pgstore

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/storesvc"
)

var (
	_ storesvc.Repository = (*PgStore)(nil)
)

type Transactor interface {
	ExtContext(ctx context.Context) sqlx.ExtContext
}

type PgStore struct {
	transactor Transactor
}

func New(
	transactor Transactor,
) *PgStore {
	return &PgStore{
		transactor: transactor,
	}
}
