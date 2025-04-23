package balance

import (
	dto "github.com/mrdolles/wallet-app/src/app/dto/balance"

	"log"

	"github.com/jmoiron/sqlx"
)

type BalanceRepository interface {
	TopUp(data *dto.TopUpReqDTO) error
	Get(data *dto.BalanceReqDTO) (*dto.BalanceRespDTO, error)
}

const (
	TopUp = `UPDATE wallets
	SET balance = balance + ?
	WHERE id = ?`

	Record = `INSERT INTO transactions (wallet_id, type, amount)
	values (?,'top-up',?)
	`

	Get = `SELECT balance from wallets 
	where id = ?
	`
)

var statement PreparedStatement

type PreparedStatement struct {
	get *sqlx.Stmt
}

type balanceRepo struct {
	Connection *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) BalanceRepository {
	repo := &balanceRepo{
		Connection: db,
	}
	InitPreparedStatement(repo)
	return repo
}

func (p *balanceRepo) Preparex(query string) *sqlx.Stmt {
	statement, err := p.Connection.Preparex(query)
	if err != nil {
		log.Fatalf("Failed to preparex query: %s. Error: %s", query, err.Error())
	}

	return statement
}

func InitPreparedStatement(m *balanceRepo) {
	statement = PreparedStatement{
		get: m.Preparex(Get),
	}
}

func (p *balanceRepo) TopUp(data *dto.TopUpReqDTO) error {

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

	_, err = tx.Exec(TopUp, data.Amount, data.WalletID)

	if err != nil {
		log.Println("Failed Query TopUp: ", err.Error())
		return err
	}

	_, err = tx.Exec(Record, data.WalletID, data.Amount)

	if err != nil {
		log.Println("Failed Query Create Record : ", err.Error())
		return err
	}

	return nil
}

func (p *balanceRepo) Get(data *dto.BalanceReqDTO) (*dto.BalanceRespDTO, error) {

	var resultData []*dto.BalanceRespDTO

	err := statement.get.Select(&resultData, data.WalletID)

	if err != nil {
		return nil, err
	}

	return resultData[0], nil
}
