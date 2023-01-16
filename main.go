package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jim-nnamdi/coldfinance/backend/admin"
	"github.com/jim-nnamdi/coldfinance/backend/connection"
	"github.com/jim-nnamdi/coldfinance/backend/content"
	"github.com/jim-nnamdi/coldfinance/backend/finance"
	"github.com/jim-nnamdi/coldfinance/backend/users"
	"go.uber.org/zap"
)

var (
	reqc   = finance.NewDataClient(&zap.Logger{})
	ticker = finance.NewStockTicker(&zap.Logger{}, reqc)
	crypto = finance.NewCryptos(&zap.Logger{}, reqc)
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

func GetAllCryptoData(w http.ResponseWriter, r *http.Request) {
	allcoins, err := crypto.GetAllCryptoData()
	if err != nil {
		logger.Debug("error fetching coins data", zap.Any("error", err))
		return
	}
	DataResponse(w, allcoins)
}

func GetLiveCryptoData(w http.ResponseWriter, r *http.Request) {
	res, err := crypto.GetLiveCryptoData()
	if err != nil {
		logger.Debug("error fetching live stats ...", zap.Any("error", err))
		return
	}
	DataResponse(w, res)
}

func ConvertCrypto(w http.ResponseWriter, r *http.Request) {
	coinfrom := r.FormValue("coinfrom")
	cointo := r.FormValue("cointo")
	amount := r.FormValue("amount")
	amtFloat, _ := strconv.Atoi(amount)
	res, err := crypto.ConvertCrypto(coinfrom, cointo, amtFloat)
	if err != nil {
		logger.Debug("error converting crypto", zap.Any("error", err))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"converted_from":    coinfrom,
		"converted_to":      cointo,
		"amount_to_convert": amtFloat,
		"amount_to_receive": res,
	})
}

func main() {
	log.Print("server running on 9900 ...")
	route := http.NewServeMux()

	// users
	route.HandleFunc("/users", users.GetAllUsers)
	route.HandleFunc("/register", users.Register)
	route.HandleFunc("/login", users.Login)

	// contents
	route.HandleFunc("/posts", content.GetAllPosts)
	route.HandleFunc("/post", content.GetPost)
	route.HandleFunc("/add/post", content.AddNewPost)
	route.HandleFunc("/posts/category", content.GetPostByCategory)

	// stocks
	route.HandleFunc("/stocks", allstocksdata)
	route.HandleFunc("/stocks/single", singleStockData)
	route.HandleFunc("/stocks/single/eod", singleStockDataEOD)
	route.HandleFunc("/stocks/splits", GetSplits)
	route.HandleFunc("/stocks/dividends", GetDividends)
	route.HandleFunc("/stocks/intraday", GetIntraday)

	// crypto
	route.HandleFunc("/coins", GetAllCryptoData)
	route.HandleFunc("/coins/live", GetLiveCryptoData)
	route.HandleFunc("/coins/convert", ConvertCrypto)

	// admin
	route.HandleFunc("/admin", admin.GetAllData)

	err := http.ListenAndServe(":9900", route)
	if err != nil {
		log.Fatal(err)
	}
}
