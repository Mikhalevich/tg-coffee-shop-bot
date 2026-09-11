package order

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type ID int

func (id ID) Int() int {
	return int(id)
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func IDFromInt(id int) ID {
	return ID(id)
}

func IDFromString(id string) (ID, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse int: %w", err)
	}

	return ID(intID), nil
}

type Order struct {
	ID               ID
	ChatID           msginfo.ChatID
	Status           Status
	VerificationCode string
	CurrencyID       currency.ID
	TotalPrice       int
	DailyPosition    int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Timeline         []StatusTime
	Products         []OrderedProduct
}

func (o *Order) IsSameChat(id msginfo.ChatID) bool {
	return o.ChatID == id
}

func (o *Order) CanCancel() bool {
	if o.Status == StatusConfirmed || o.Status == StatusWaitingPayment {
		return true
	}

	return false
}

func (o *Order) InQueue() bool {
	if o.Status == StatusConfirmed || o.Status == StatusInProgress {
		return true
	}

	return false
}

func (o *Order) ProductIDs() []product.ProductID {
	ids := make([]product.ProductID, 0, len(o.Products))
	for _, v := range o.Products {
		ids = append(ids, v.ProductID)
	}

	return ids
}

type StatusTime struct {
	Status Status
	Time   time.Time
}

type OrderedProduct struct {
	ProductID  product.ProductID
	CategoryID product.CategoryID // available only for cart products.
	Count      int
	Price      int
}

type CreateOrderInfo struct {
	ChatID              msginfo.ChatID
	Status              Status
	StatusOperationTime time.Time
	VerificationCode    string
	TotalPrice          int
	Products            []OrderedProduct
	CurrencyID          currency.ID
}
