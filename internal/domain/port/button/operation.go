package button

type Operation string

const (
	OperationOpenURL               Operation = "OpenURL"
	OperationOpenBox               Operation = "OpenBox"
	OperationCardPage              Operation = "CardPage"
	OperationCardTotal             Operation = "CardTotal"
	OperationAbstractDuplicatesAll Operation = "AbstractDuplicatesAll"
	OperationAbstractCard          Operation = "AbstractCard"
	OperationBuyBox                Operation = "BuyBox"
	OperationShop                  Operation = "Shop"
	OperationGetCommonBox          Operation = "GetCommonBox"
	OperationLike                  Operation = "Like"
	OperationBoxReadyToOpen        Operation = "BoxReadyToOpen"
)

func (o Operation) String() string {
	return string(o)
}
