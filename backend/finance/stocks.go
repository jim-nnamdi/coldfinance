package finance

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type StockTicker interface {
	GetAllTickers() (allStockTickers, error)
	GetCompanyTicker(companyTicker string) (*allTickers, error)
	GetCompanyEOD(symbol string) (*ParentStockEOD, error)
	GetCompanySplits(symbol string) (*Split, error)
	GetCompanyIntraday(symbol string) (*Intraday, error)
}

var _ StockTicker = &stockTickers{}

type stockTickers struct {
	logger      *zap.Logger
	stockClient RequestClient
}

type allStockTickers struct {
	// Pagination StockPagination `json:"pagination"`
	Data []allTickers `json:"data"`
}

type StockPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
	Total  int `json:"total"`
}

type allTickers struct {
	Name          string        `json:"name"`
	Symbol        string        `json:"symbol"`
	HasIntraDay   bool          `json:"has_intraday"`
	HasEOD        bool          `json:"has_eod"`
	Country       string        `json:"country"`
	StockExchange StockExchange `json:"stock_exchange"`
}

type StockExchange struct {
	Name        string `json:"name"`
	Acronym     string `json:"acronym"`
	MIC         string `json:"mic"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
	Website     string `json:"website"`
}

type ParentStockEOD struct {
	// Pagination StockPagination `json:"pagination"`
	Data StockEOD `json:"data"`
}

type StockEOD struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	HasIntraDay bool   `json:"has_intraday"`
	HasEOD      bool   `json:"has_eod"`
	Country     string `json:"country"`
	EOD         []EOD  `json:"eod"`
}

type EOD struct {
	Open        float64 `json:"open"`
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Close       float64 `json:"close"`
	Volume      float64 `json:"volume"`
	AdjHigh     float64 `json:"adj_high"`
	AdjLow      float64 `json:"adj_low"`
	AdjClose    float64 `json:"adj_close"`
	AdjOpen     float64 `json:"adj_open"`
	AdjVolume   float64 `json:"adj_volume"`
	SplitFactor float64 `json:"split_factor"`
	Dividend    float64 `json:"dividend"`
	Symbol      string  `json:"symbol"`
	Exchange    string  `json:"exchange"`
	Date        string  `json:"date"`
}

type Split struct {
	Data []struct {
		Date        string  `json:"date"`
		SplitFactor float64 `json:"split_factor"`
		Symbol      string  `json:"symbol"`
	} `json:"data"`
}

type Dividend struct {
	Data []struct {
		Date     string  `json:"date"`
		Dividend float64 `json:"dividend"`
		Symbol   string  `json:"symbol"`
	} `json:"data"`
}

type Intraday struct {
	Data struct {
		Name        string `json:"name"`
		Symbol      string `json:"symbol"`
		Country     string `json:"country"`
		HasIntraday bool   `json:"has_intraday"`
		HasEOD      bool   `json:"has_eod"`
		Intraday    []struct {
			Open     float64 `json:"open"`
			High     float64 `json:"high"`
			Low      float64 `json:"low"`
			Last     string  `json:"last"`
			Close    string  `json:"close"`
			Volume   string  `json:"volume"`
			Date     string  `json:"date"`
			Symbol   string  `json:"symbol"`
			Exchange string  `json:"exchange"`
		} `json:"intraday"`
	} `json:"data"`
}

func NewStockTicker(logger *zap.Logger, stockclient RequestClient) *stockTickers {
	return &stockTickers{
		logger:      logger,
		stockClient: stockclient,
	}
}

func (s *stockTickers) GetAllTickers() (allStockTickers, error) {
	req, err := http.NewRequest(http.MethodGet, "http://api.marketstack.com/v1/tickers?access_key=214e04d7d20abb7cbd47894b00c2c334", nil)
	if err != nil {
		s.logger.Debug("error fetching stock tickers", zap.Any("error", err.Error()))
		return allStockTickers{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		s.logger.Debug("cannot process response", zap.Any("error", err.Error()))
		return allStockTickers{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Debug("failed to process response byte", zap.Any("error", err.Error()))
		return allStockTickers{}, err
	}

	var val allStockTickers
	data := json.Unmarshal(body, &val)
	if data != nil {
		s.logger.Debug("error mapping data to struct", zap.Any("error", data))
		return allStockTickers{}, err
	}
	return val, nil
}

func (s *stockTickers) GetCompanyTicker(symbol string) (*allTickers, error) {
	err := godotenv.Load()
	if err != nil {
		s.logger.Error("error loading env file", zap.Error(err))
		return nil, err
	}
	accessKey := os.Getenv("MARKETSTACK")
	req, err := http.NewRequest(http.MethodGet, "http://api.marketstack.com/v1/tickers/"+symbol+"?access_key="+accessKey, nil)
	if err != nil {
		s.logger.Debug("error fetching stock tickers", zap.Any("error", err.Error()))
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		s.logger.Debug("error", zap.String("error", err.Error()))
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Debug("error", zap.String("error", err.Error()))
		return nil, err
	}
	var val allTickers
	datast := json.Unmarshal(body, &val)
	if datast != nil {
		s.logger.Debug("error unmarshaling to struct", zap.Any("error", datast))
		return nil, err
	}
	return &val, nil
}

func (s *stockTickers) GetCompanyEOD(symbol string) (*ParentStockEOD, error) {
	err := godotenv.Load()
	if err != nil {
		s.logger.Error("error loading env file", zap.Error(err))
		return nil, err
	}
	accessKey := os.Getenv("MARKETSTACK")
	req, err := s.stockClient.MakeGetRequest(http.MethodGet, "http://api.marketstack.com/v1/tickers/"+symbol+"/eod?access_key="+accessKey)
	if err != nil {
		s.logger.Debug("error making request", zap.Any("error", err))
		return nil, err
	}
	var val ParentStockEOD
	datast := json.Unmarshal(req, &val)
	if datast != nil {
		s.logger.Debug("error unmarshalling data", zap.Any("error", datast))
		return nil, datast
	}
	return &val, nil
}

func (s *stockTickers) GetCompanySplits(symbol string) (*Split, error) {
	err := godotenv.Load()
	if err != nil {
		s.logger.Error("error loading env file", zap.Error(err))
		return nil, err
	}
	accessKey := os.Getenv("MARKETSTACK")
	req, err := s.stockClient.MakeGetRequest(http.MethodGet, "http://api.marketstack.com/v1/tickers/"+symbol+"/splits?access_key="+accessKey)
	if err != nil {
		s.logger.Debug("error making request", zap.Any("error", err))
		return nil, err
	}
	var val Split
	datast := json.Unmarshal(req, &val)
	if datast != nil {
		s.logger.Debug("error unmarshalling data", zap.Any("error", datast))
		return nil, datast
	}
	return &val, nil
}

func (s *stockTickers) GetCompanyDividends(symbol string) (*Dividend, error) {
	err := godotenv.Load()
	if err != nil {
		s.logger.Error("error loading env file", zap.Error(err))
		return nil, err
	}
	accessKey := os.Getenv("MARKETSTACK")
	req, err := s.stockClient.MakeGetRequest(http.MethodGet, "http://api.marketstack.com/v1/tickers/"+symbol+"/dividends?access_key="+accessKey)
	if err != nil {
		s.logger.Debug("error making request", zap.Any("error", err))
		return nil, err
	}
	log.Print("dividends", string(req))
	var val Dividend
	datast := json.Unmarshal(req, &val)
	if datast != nil {
		s.logger.Debug("error unmarshalling data", zap.Any("error", datast))
		return nil, datast
	}
	return &val, nil
}

func (s *stockTickers) GetCompanyEODLatest(symbol string) (*ParentStockEOD, error) {
	err := godotenv.Load()
	if err != nil {
		s.logger.Error("error loading env file", zap.Error(err))
		return nil, err
	}
	accessKey := os.Getenv("MARKETSTACK")
	req, err := s.stockClient.MakeGetRequest(http.MethodGet, "http://api.marketstack.com/v1/tickers/"+symbol+"/eod/latest?access_key="+accessKey)
	if err != nil {
		s.logger.Debug("error making request", zap.Any("error", err))
		return nil, err
	}
	var val ParentStockEOD
	datast := json.Unmarshal(req, &val)
	if datast != nil {
		s.logger.Debug("error unmarshalling data", zap.Any("error", datast))
		return nil, datast
	}
	return &val, nil
}

func (s *stockTickers) GetCompanyIntraday(symbol string) (*Intraday, error) {
	err := godotenv.Load()
	if err != nil {
		s.logger.Error("error loading env file", zap.Error(err))
		return nil, err
	}
	accessKey := os.Getenv("MARKETSTACK")
	req, err := s.stockClient.MakeGetRequest(http.MethodGet, "http://api.marketstack.com/v1/tickers/"+symbol+"/intraday?access_key="+accessKey)
	if err != nil {
		s.logger.Debug("error making request", zap.Any("error", err))
		return nil, err
	}
	var val Intraday
	_ = json.Unmarshal(req, &val)
	return &val, nil
}

func (s *stockTickers) GetCompanyIntradayLatest(symbol string) (*ParentStockEOD, error) {
	err := godotenv.Load()
	if err != nil {
		s.logger.Error("error loading env file", zap.Error(err))
		return nil, err
	}
	accessKey := os.Getenv("MARKETSTACK")
	req, err := s.stockClient.MakeGetRequest(http.MethodGet, "http://api.marketstack.com/v1/tickers/"+symbol+"/intraday/latest?access_key="+accessKey)
	if err != nil {
		s.logger.Debug("error making request", zap.Any("error", err))
		return nil, err
	}
	var val ParentStockEOD
	datast := json.Unmarshal(req, &val)
	if datast != nil {
		s.logger.Debug("error unmarshalling data", zap.Any("error", datast))
		return nil, datast
	}
	return &val, nil
}
