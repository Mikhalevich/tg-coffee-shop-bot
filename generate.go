package generate

//go:generate go tool mockgen -source=internal/domain/customer/orderaction/order_action.go -destination=internal/domain/customer/orderaction/order_action_mock.go -package=orderaction

//go:generate go tool mockgen -source=internal/adapter/repository/postgres/transaction/transaction.go -destination=internal/adapter/repository/postgres/transaction/transaction_mock.go -package=transaction
