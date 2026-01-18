package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/cartprocessing"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderaction"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderhistory"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderpayment"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/manager/orderprocessing"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/outboxprocessor"
)

var (
	_ cartprocessing.Repository = (*Postgres)(nil)
	_ cartprocessing.StoreInfo  = (*Postgres)(nil)

	_ orderaction.Repository = (*Postgres)(nil)

	_ orderpayment.Repository    = (*Postgres)(nil)
	_ orderpayment.MessageSender = (*Postgres)(nil)

	_ orderpayment.StoreInfo = (*Postgres)(nil)

	_ orderhistory.CurrencyProvider = (*Postgres)(nil)

	_ orderprocessing.Repository                  = (*Postgres)(nil)
	_ orderprocessing.CustomerOutboxMessageSender = (*Postgres)(nil)

	_ outboxprocessor.Repository = (*Postgres)(nil)
)

type Driver interface {
	IsConstraintError(err error, constraint string) bool
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
	ExtContext(ctx context.Context) sqlx.ExtContext
}

type Postgres struct {
	driver     Driver
	transactor Transactor
}

func New(
	driver Driver,
	transactor Transactor,
) *Postgres {
	return &Postgres{
		driver:     driver,
		transactor: transactor,
	}
}

func (p *Postgres) Transactor() Transactor {
	return p.transactor
}
