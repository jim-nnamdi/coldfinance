package admin

import (
	"encoding/json"
	"net/http"

	"github.com/jim-nnamdi/coldfinance/backend/connection"
	"go.uber.org/zap"
)

var (
	conn           = connection.Dbconn()
	coldfinancelog = connection.Coldfinancelog()
)

func GetPostsCount() (int, error) {
	postcount, err := conn.Query("select count(*) from posts")
	if err != nil {
		coldfinancelog.Debug("could not fetch data from db", zap.Any("data", postcount))
		return 0, err
	}
	defer postcount.Close()
	var count int
	for postcount.Next() {
		if err := postcount.Scan(&count); err != nil {
			coldfinancelog.Debug("could not fetch count", zap.Any("error", err))
			return 0, err
		}
	}
	return count, nil
}

func GetUsersCount() (int, error) {
	usercount, err := conn.Query("select count(*) from users")
	if err != nil {
		coldfinancelog.Debug("could not fetch data from db", zap.Any("data", usercount))
		return 0, err
	}
	defer usercount.Close()
	var count int
	for usercount.Next() {
		if err := usercount.Scan(&count); err != nil {
			coldfinancelog.Debug("could not fetch count", zap.Any("error", err))
			return 0, err
		}
	}
	return count, nil
}

func GetAllData(w http.ResponseWriter, r *http.Request) {
	ucount, _ := GetUsersCount()
	pcount, _ := GetPostsCount()
	retdata := map[string]interface{}{}
	retdata["usercount"] = ucount
	retdata["postcount"] = pcount

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(retdata)
}
