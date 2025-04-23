package transaction

import (
	dto "github.com/mrdolles/wallet-app/src/app/dto/transaction"
	"github.com/mrdolles/wallet-app/src/infra/helper"

	"log"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	Transfer(data *dto.TransferReqDTO) error
	GetTopTen(walletID int64) ([]*dto.GetTopTenRespDTO, error)
	GetOverallTopTransactions() ([]*dto.GetOverallRespDTO, error)
}

const (
	GetTopTenByUser = `SELECT 
	u.username,
    t.amount
	FROM 
	    transactions t
	INNER JOIN 
	    wallets w ON t.wallet_id = w.id
	INNER JOIN 
	    users u ON w.user_id = u.id
	WHERE 
	    t.type = 'transfer' AND w.id = ?
	ORDER BY 
	    ABS(t.amount) DESC
	LIMIT 10;
	`

	GetOverallTopTransactions = `
	SELECT 
    u.username,
    SUM(ABS(t.amount)) AS transacted_value
	FROM 
	    transactions t
	INNER JOIN 
	    wallets w ON t.wallet_id = w.id
	INNER JOIN 
	    users u ON w.user_id = u.id
	WHERE 
	    t.type = 'transfer'
	    AND t.amount < 0
	GROUP BY 
	    u.username
	ORDER BY 
	    transacted_value DESC
	LIMIT 10;
	`
)

var statement PreparedStatement

type PreparedStatement struct {
	getTopTen  *sqlx.Stmt
	getOverall *sqlx.Stmt
}

type transactionRepo struct {
	Connection *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	repo := &transactionRepo{
		Connection: db,
	}
	InitPreparedStatement(repo)
	return repo
}

func (p *transactionRepo) Preparex(query string) *sqlx.Stmt {
	statement, err := p.Connection.Preparex(query)
	if err != nil {
		log.Fatalf("Failed to preparex query: %s. Error: %s", query, err.Error())
	}

	return statement
}

func InitPreparedStatement(m *transactionRepo) {
	statement = PreparedStatement{
		getTopTen:  m.Preparex(GetTopTenByUser),
		getOverall: m.Preparex(GetOverallTopTransactions),
	}
}

func (p *transactionRepo) Transfer(data *dto.TransferReqDTO) error {

	tx, err := p.Connection.Beginx()
	if err != nil {
		log.Println("Failed Begin Tx TopUp  : ", err.Error())
		return err
	}
	defer func(tx *sqlx.Tx) {
		if err != nil {
			tx.Rollback()
			log.Println("Rolling back transaction due to:", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Println("Failed to commit transaction:", err.Error())
			}
		}
	}(tx)

	resultTransfer, err := tx.Exec("UPDATE wallets SET balance = balance - ? WHERE id = ? AND balance >= ?", data.Amount, data.WalletIDSender, data.Amount)
	if err != nil {
		log.Println("Failed Query Transfer: ", err.Error())
		return err
	}

	row, _ := resultTransfer.RowsAffected()
	if row < 1 {
		log.Println("Failed Query Transfer: ", helper.ErrInsufficientBalance)
		err = helper.ErrInsufficientBalance
		return err
	}

	_, err = tx.Exec("INSERT INTO transactions (wallet_id, type, amount) VALUES (?, 'transfer', -CAST(? AS DECIMAL))", data.WalletIDSender, data.Amount)

	if err != nil {
		log.Println("Failed Query Create Transaction Transfer : ", err.Error())
		return err
	}

	resultReceive, err := tx.Exec("UPDATE wallets w JOIN ( SELECT w.id FROM wallets w INNER JOIN users u ON u.id = w.user_id WHERE u.username = ? LIMIT 1 ) AS sub ON w.id = sub.id SET w.balance = w.balance + ?", data.ToUserName, data.Amount)

	if err != nil {
		log.Println("Failed Query Create Receive : ", err.Error())
		return err
	}

	rowReceive, _ := resultReceive.RowsAffected()
	if rowReceive < 1 {
		log.Println("Failed Query Receive: ", helper.ErrUserNotFound)
		err = helper.ErrUserNotFound
		return err
	}

	_, err = tx.Exec("INSERT INTO transactions (wallet_id, type, amount) VALUES ((SELECT w.id FROM wallets w INNER JOIN users u ON u.id = w.user_id WHERE u.username = ? LIMIT 1), 'transfer', ?)", data.ToUserName, data.Amount)

	if err != nil {
		log.Println("Failed Query Create Transaction Receive : ", err.Error())
		return err
	}

	return nil
}

func (p *transactionRepo) GetTopTen(walletID int64) ([]*dto.GetTopTenRespDTO, error) {

	var resultData []*dto.GetTopTenRespDTO

	err := statement.getTopTen.Select(&resultData, walletID)

	if err != nil {
		return nil, err
	}

	return resultData, nil
}

func (p *transactionRepo) GetOverallTopTransactions() ([]*dto.GetOverallRespDTO, error) {

	var resultData []*dto.GetOverallRespDTO

	err := statement.getOverall.Select(&resultData)

	if err != nil {
		return nil, err
	}

	return resultData, nil
}
