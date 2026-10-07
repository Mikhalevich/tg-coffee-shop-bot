package outboxmsg

type AnswerPayment struct {
	ID        int
	PaymentID string
	OK        bool
	ErrorMsg  string
}
