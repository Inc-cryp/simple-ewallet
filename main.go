package main

import (
	"context"
	"database/sql"

	usecases "github.com/mrdolles/wallet-app/src/app/usecases"

	"github.com/mrdolles/wallet-app/src/infra/config"

	mysql "github.com/mrdolles/wallet-app/src/infra/persistence/mysql"
	"github.com/mrdolles/wallet-app/src/infra/persistence/redis"

	balanceRepo "github.com/mrdolles/wallet-app/src/infra/persistence/mysql/balance"
	transactionRepo "github.com/mrdolles/wallet-app/src/infra/persistence/mysql/transaction"
	userRepo "github.com/mrdolles/wallet-app/src/infra/persistence/mysql/user"

	redisServe "github.com/mrdolles/wallet-app/src/infra/persistence/redis/service"
	"github.com/mrdolles/wallet-app/src/interface/rest"

	ms_log "github.com/mrdolles/wallet-app/src/infra/log"

	balanceUC "github.com/mrdolles/wallet-app/src/app/usecases/balance"
	transactionUC "github.com/mrdolles/wallet-app/src/app/usecases/transaction"
	userUC "github.com/mrdolles/wallet-app/src/app/usecases/user"

	_ "github.com/joho/godotenv/autoload"

	"github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	conf := config.Make()

	isProd := false
	if conf.App.Environment == "PRODUCTION" {
		isProd = true
	}

	m := make(map[string]interface{})
	m["env"] = conf.App.Environment
	m["service"] = conf.App.Name
	logger := ms_log.NewLogInstance(
		ms_log.LogName(conf.Log.Name),
		ms_log.IsProduction(isProd),
		ms_log.LogAdditionalFields(m))

	mysqldb, _ := mysql.New(conf.SqlDb, logger)
	redisClient, _ := redis.NewRedisClient(conf.Redis, logger)

	redisServe := redisServe.NewServRedis(redisClient)

	defer func(l *logrus.Logger, sqlDB *sql.DB, dbName string) {
		err := sqlDB.Close()
		if err != nil {
			l.Errorf("error closing sql database %s: %s", dbName, err)
		} else {
			l.Printf("sql database %s successfuly closed.", dbName)
		}
	}(logger, mysqldb.Conn.DB, mysqldb.Conn.DriverName())

	userRepository := userRepo.NewUserRepository(mysqldb.Conn)
	balanceRepository := balanceRepo.NewBalanceRepository(mysqldb.Conn)
	transactionRepository := transactionRepo.NewTransactionRepository(mysqldb.Conn)

	httpServer, err := rest.New(
		conf.Http,
		isProd,
		logger,
		usecases.AllUseCases{
			UserUC:        userUC.NewUserUseCase(userRepository),
			BalanceUC:     balanceUC.NewBalanceUseCase(balanceRepository),
			TransactionUC: transactionUC.NewTransactionUseCase(transactionRepository, redisServe),
		},
	)
	if err != nil {
		panic(err)
	}
	httpServer.Start(ctx)

}
