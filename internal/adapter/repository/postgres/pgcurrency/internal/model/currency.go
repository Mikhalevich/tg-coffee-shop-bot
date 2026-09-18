package model

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
)

type Currency struct {
	ID         int    `db:"id"`
	Code       string `db:"code"`
	Exp        int    `db:"exp"`
	DecimalSep string `db:"decimal_sep"`
	MinAmount  int    `db:"min_amount"`
	MaxAmount  int    `db:"max_amount"`
	IsEnabled  bool   `db:"is_enabled"`
}

func (c Currency) ToDom() currency.Currency {
	return currency.Currency{
		ID:         currency.IDFromInt(c.ID),
		Code:       c.Code,
		Exp:        c.Exp,
		DecimalSep: c.DecimalSep,
		MinAmount:  c.MinAmount,
		MaxAmount:  c.MaxAmount,
		IsEnabled:  c.IsEnabled,
	}
}
