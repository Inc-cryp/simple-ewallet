package usecases

import (
	balanceUC "github.com/mrdolles/wallet-app/src/app/usecases/balance"
	transactionUC "github.com/mrdolles/wallet-app/src/app/usecases/transaction"
	userUC "github.com/mrdolles/wallet-app/src/app/usecases/user"
)

type AllUseCases struct {
	UserUC        userUC.UserUCInterface
	BalanceUC     balanceUC.BalanceUCInterface
	TransactionUC transactionUC.TransactionUCInterface
}
