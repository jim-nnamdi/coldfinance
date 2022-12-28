package connection

import (
	"database/sql"
	"log"

	"go.uber.org/zap"
)

func Coldfinancelog() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	return sugar
}

func Dbconn() *sql.DB {
	db, err := sql.Open("mysql", "root:M@etroboomin50@tcp(localhost:3306)/coldfinance")
	if err != nil {
		log.Print(err)
		panic(err)
	}
	return db
}
