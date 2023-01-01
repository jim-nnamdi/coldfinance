package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jim-nnamdi/coldfinance/backend/admin"
	"github.com/jim-nnamdi/coldfinance/backend/connection"
	"github.com/jim-nnamdi/coldfinance/backend/content"
	"github.com/jim-nnamdi/coldfinance/backend/finance"
	"github.com/jim-nnamdi/coldfinance/backend/users"
	"go.uber.org/zap"
)

var (
	stockc = finance.NewStockClient(&zap.Logger{})
	ticker = finance.NewStockTicker(&zap.Logger{}, stockc)
	logger = zap.NewNop()
)

func DataResponse(w http.ResponseWriter, v any) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func allstocksdata(w http.ResponseWriter, r *http.Request) {
	stocks, err := ticker.GetAllTickers()
	if err != nil {
		connection.Coldfinancelog().Debug("error", zap.Any("error", err))
		return
	}
	DataResponse(w, stocks)
}

func singleStockData(w http.ResponseWriter, r *http.Request) {
	sym := r.FormValue("symbol")
	res, err := ticker.GetCompanyTicker(sym)
	if err != nil {
		logger.Debug("cannot process single stock data", zap.Any("error", err))
		return
	}
	DataResponse(w, res)
}

func singleStockDataEOD(w http.ResponseWriter, r *http.Request) {
	sym := r.FormValue("symbol")
	res, err := ticker.GetCompanyEOD(sym)
	if err != nil {
		logger.Debug("cannot fetch EOD for symbol", zap.Any("error", err))
		return
	}
	DataResponse(w, res)
}

func GetSplits(w http.ResponseWriter, r *http.Request) {
	sym := r.FormValue("symbol")
	res, err := ticker.GetCompanySplits(sym)
	if err != nil {
		logger.Debug("cannot fetch EOD for symbol", zap.Any("error", err))
		return
	}
	DataResponse(w, res)
}

func GetDividends(w http.ResponseWriter, r *http.Request) {
	sym := r.FormValue("symbol")
	res, err := ticker.GetCompanyDividends(sym)
	if err != nil {
		logger.Debug("cannot fetch EOD for symbol", zap.Any("error", err))
		return
	}
	DataResponse(w, res)
}

func GetIntraday(w http.ResponseWriter, r *http.Request) {
	sym := r.FormValue("symbol")
	res, err := ticker.GetCompanyIntraday(sym)
	if err != nil {
		logger.Debug("cannot fetch EOD for symbol", zap.Any("error", err))
		return
	}
	DataResponse(w, res)
}

func main() {
	log.Print("server running on 9900 ...")
	r := http.NewServeMux()
	r.HandleFunc("/users", users.GetAllUsers)
	r.HandleFunc("/register", users.Register)
	r.HandleFunc("/login", users.Login)

	r.HandleFunc("/posts", content.GetAllPosts)
	r.HandleFunc("/post", content.GetPost)
	r.HandleFunc("/add/post", content.AddNewPost)

	r.HandleFunc("/admin", admin.GetAllData)

	// stocks
	r.HandleFunc("/stocks", allstocksdata)
	r.HandleFunc("/stocks/single", singleStockData)
	r.HandleFunc("/stocks/single/eod", singleStockDataEOD)
	r.HandleFunc("/stocks/splits", GetSplits)
	r.HandleFunc("/stocks/dividends", GetDividends)
	r.HandleFunc("/stocks/intraday", GetIntraday)
	err := http.ListenAndServe(":9900", r)
	if err != nil {
		log.Fatal(err)
	}
}
