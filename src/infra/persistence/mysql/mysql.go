package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/mrdolles/wallet-app/src/infra/config"
	"github.com/sirupsen/logrus"
)

type MysqlDb struct {
	Conn *sqlx.DB
}

func New(conf config.SqlDbConf, logger *logrus.Logger) (MysqlDb, error) {
	db := MysqlDb{}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name,
	)
	if conf.Password == "" {
		dsn = fmt.Sprintf(
			"%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			conf.Username,
			conf.Host,
			conf.Port,
			conf.Name,
		)
	}

	conn, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect to database")
	}

	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}

	logger.Printf("sql database connection %s seccess", db.Conn.Driver())
	return db, nil
}
