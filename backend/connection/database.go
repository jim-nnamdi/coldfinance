package connection

import (
	"database/sql"
	"log"
)

func Dbconn() *sql.DB {
	db, err := sql.Open("mysql", "root:M@etroboomin50@tcp(localhost:3306)/coldfinance")
	if err != nil {
		log.Print(err)
		panic(err)
	}
	return db
}
