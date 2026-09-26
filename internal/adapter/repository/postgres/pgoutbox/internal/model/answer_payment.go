package model

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/outboxmsg"
)

type AnswerPayment struct {
	ID        int    `db:"id"`
	PaymentID string `db:"payment_id"`
	OK        bool   `db:"ok"`
	ErrorMsg  string `db:"error_msg"`
}

func (ap AnswerPayment) ToDom() outboxmsg.AnswerPayment {
	return outboxmsg.AnswerPayment{
		ID:        ap.ID,
		PaymentID: ap.PaymentID,
		OK:        ap.OK,
		ErrorMsg:  ap.ErrorMsg,
	}
}

func ToDomAnswerPayments(dbPayments []AnswerPayment) []outboxmsg.AnswerPayment {
	if len(dbPayments) == 0 {
		return nil
	}

	payments := make([]outboxmsg.AnswerPayment, 0, len(dbPayments))

	for _, p := range dbPayments {
		payments = append(payments, p.ToDom())
	}

	return payments
}
