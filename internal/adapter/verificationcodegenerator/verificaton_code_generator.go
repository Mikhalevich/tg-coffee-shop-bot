package verificationcodegenerator

import (
	"fmt"
	"math/rand"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/payment"
)

const (
	maxCodeValue = 1000
)

var _ payment.CodeGeneratorService = (*VerificationCodeGenerator)(nil)

type VerificationCodeGenerator struct {
}

func New() *VerificationCodeGenerator {
	return &VerificationCodeGenerator{}
}

func (g *VerificationCodeGenerator) Generate() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(maxCodeValue))
}
