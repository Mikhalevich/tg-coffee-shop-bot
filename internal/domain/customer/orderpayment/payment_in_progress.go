package orderpayment

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

func (p *OrderPayment) PaymentInProgress(
	ctx context.Context,
	paymentID string,
	orderID order.ID,
	currency string,
	totalAmount int,
) error {
	storeInfo, err := p.storeInfoByID(ctx, p.storeID)
	if err != nil {
		return fmt.Errorf("check for active: %w", err)
	}

	if !storeInfo.IsActive {
		if err := p.sender.AnswerOrderPayment(ctx, paymentID, false, storeInfo.ClosedStoreMessage); err != nil {
			return fmt.Errorf("answer payment for store is closed: %w", err)
		}

		return nil
	}

	if err := p.transactor.Transaction(ctx,
		func(ctx context.Context) error {
			res, err := p.setOrderInProgress(ctx, orderID, totalAmount)
			if err != nil {
				return fmt.Errorf("set order in progress: %w", err)
			}

			if err := p.sender.AnswerOrderPayment(ctx, paymentID, res.OK, res.ErrorMsg); err != nil {
				return fmt.Errorf("answer order payment: %w", err)
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

type answerOrderPaymentResult struct {
	OK       bool
	ErrorMsg string
}

func (p *OrderPayment) setOrderInProgress(
	ctx context.Context,
	orderID order.ID,
	totalAmount int,
) (*answerOrderPaymentResult, error) {
	ord, err := p.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if p.repository.IsNotFoundError(err) {
			return &answerOrderPaymentResult{
				OK:       false,
				ErrorMsg: message.OrderNotExists(),
			}, nil
		}

		return nil, fmt.Errorf("get order by id: %w", err)
	}

	if ord.Status != order.StatusWaitingPayment {
		return &answerOrderPaymentResult{
			OK:       false,
			ErrorMsg: message.OrderStatus(ord.Status),
		}, nil
	}

	if ord.TotalPrice != totalAmount {
		return &answerOrderPaymentResult{
			OK:       false,
			ErrorMsg: message.OrderTotalPriceIncorrect(),
		}, nil
	}

	if _, err := p.repository.UpdateOrderStatus(
		ctx,
		orderID,
		p.timeProvider.Now(),
		order.StatusPaymentInProgress,
		order.StatusWaitingPayment,
	); err != nil {
		return nil, fmt.Errorf("update order status: %w", err)
	}

	return &answerOrderPaymentResult{
		OK: true,
	}, nil
}

type storeInfo struct {
	Store              *store.Store
	IsActive           bool
	ClosedStoreMessage string
}

func (p *OrderPayment) storeInfoByID(ctx context.Context, storeID store.ID) (*storeInfo, error) {
	stor, err := p.storeInfo.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("get store by id: %w", err)
	}

	currentTime := p.timeProvider.Now()

	nextWorkingTime, isActive := stor.Schedule.NextWorkingTime(currentTime)
	if !isActive {
		return &storeInfo{
			Store:              stor,
			IsActive:           false,
			ClosedStoreMessage: message.StoreClosed(currentTime, nextWorkingTime),
		}, nil
	}

	return &storeInfo{
		Store:    stor,
		IsActive: true,
	}, nil
}
