package model

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type OrderProduct struct {
	OrderID   int `db:"order_id"`
	ProductID int `db:"product_id"`
	Count     int `db:"count"`
	Price     int `db:"price"`
}

func ToDBOrderProducts(id order.ID, portProducts []order.OrderedProduct) []OrderProduct {
	dbProducts := make([]OrderProduct, 0, len(portProducts))

	for _, v := range portProducts {
		dbProducts = append(dbProducts, OrderProduct{
			OrderID:   id.Int(),
			ProductID: v.ProductID.Int(),
			Count:     v.Count,
			Price:     v.Price,
		})
	}

	return dbProducts
}

func toDomOrderedProducts(dbProducts []OrderProduct) []order.OrderedProduct {
	domProducts := make([]order.OrderedProduct, 0, len(dbProducts))

	for _, v := range dbProducts {
		domProducts = append(domProducts, order.OrderedProduct{
			ProductID: product.ProductIDFromInt(v.ProductID),
			Count:     v.Count,
			Price:     v.Price,
		})
	}

	return domProducts
}
