package model

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/outboxprocessor"
)

type OutboxAnswerPayment struct {
	ID        int    `db:"id"`
	PaymentID string `db:"payment_id"`
	OK        bool   `db:"ok"`
	ErrorMsg  string `db:"error_msg"`
}

func ToOutboxAnswerPayment(dbPayment OutboxAnswerPayment) outboxprocessor.OutboxAnswerPayment {
	return outboxprocessor.OutboxAnswerPayment{
		ID:        dbPayment.ID,
		PaymentID: dbPayment.PaymentID,
		OK:        dbPayment.OK,
		ErrorMsg:  dbPayment.ErrorMsg,
	}
}

func ToOutboxAnswerPayments(dbPayments []OutboxAnswerPayment) []outboxprocessor.OutboxAnswerPayment {
	if len(dbPayments) == 0 {
		return nil
	}

	payments := make([]outboxprocessor.OutboxAnswerPayment, 0, len(dbPayments))

	for _, p := range dbPayments {
		payments = append(payments, ToOutboxAnswerPayment(p))
	}

	return payments
}
