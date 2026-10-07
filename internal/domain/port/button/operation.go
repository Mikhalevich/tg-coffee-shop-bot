package button

type Operation string

const (
	OperationOpenURL Operation = "OpenURL"

	OperationOrderCancel Operation = "OrderCancel"

	OperationOrderHistoryByPage      Operation = "OperationOrderHistoryByPage"
	OperationOrderHistoryByPageFirst Operation = "OperationOrderHistoryByPageFirst"
	OperationOrderHistoryByPageLast  Operation = "OperationOrderHistoryByPageLast"

	OperationCartCancel               Operation = "CartCancel"
	OperationCartConfirm              Operation = "CartConfirm"
	OperationCartViewCategoryProducts Operation = "CartViewCategoryProducts"
	OperationCartViewCategories       Operation = "CartViewCategories"
	OperationCartAddProduct           Operation = "CartAddProduct"
)

func (o Operation) String() string {
	return string(o)
}
