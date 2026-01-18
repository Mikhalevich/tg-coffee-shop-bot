package orderaction_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderaction"
)

type OrderActionSuite struct {
	*suite.Suite

	ctrl *gomock.Controller

	sender     *orderaction.MockMessageSender
	repository *orderaction.MockRepository

	orderAction *orderaction.OrderAction
}

func TestProcessorSuit(t *testing.T) {
	t.Parallel()
	suite.Run(t, &OrderActionSuite{
		Suite: new(suite.Suite),
	})
}

func (s *OrderActionSuite) SetupSuite() {
	s.ctrl = gomock.NewController(s.T())

	s.sender = orderaction.NewMockMessageSender(s.ctrl)
	s.repository = orderaction.NewMockRepository(s.ctrl)

	s.orderAction = orderaction.New(s.sender, s.repository, nil)
}

func (s *OrderActionSuite) TearDownSuite() {
}

func (s *OrderActionSuite) TearDownTest() {
}

func (s *OrderActionSuite) TearDownSubTest() {
}
