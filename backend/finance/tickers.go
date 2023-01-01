package finance

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.uber.org/zap"
)

type StockTicker interface {
	GetAllTickers() (allStockTickers, error)
	GetCompanyTicker(companyTicker string) (allTickers, error)
	GetCompanyEOD(symbol string) (*ParentStockEOD, error)
}

var _ StockTicker = &stockTickers{}

type stockTickers struct {
	logger      *zap.Logger
	stockClient Stockclient
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

func NewStockTicker(logger *zap.Logger, stockclient Stockclient) *stockTickers {
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

func (s *stockTickers) GetCompanyTicker(symbol string) (allTickers, error) {
	req, err := http.NewRequest(http.MethodGet, "http://api.marketstack.com/v1/tickers/"+symbol+"?access_key=214e04d7d20abb7cbd47894b00c2c334", nil)
	if err != nil {
		s.logger.Debug("error fetching stock tickers", zap.Any("error", err.Error()))
		return allTickers{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		s.logger.Debug("error", zap.String("error", err.Error()))
		return allTickers{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Debug("error", zap.String("error", err.Error()))
		return allTickers{}, err
	}
	var val allTickers
	datast := json.Unmarshal(body, &val)
	if datast != nil {
		s.logger.Debug("error unmarshaling to struct", zap.Any("error", datast))
		return allTickers{}, err
	}
	return val, nil
}

func (s *stockTickers) GetCompanyEOD(symbol string) (*ParentStockEOD, error) {
	req, err := s.stockClient.MakeGetRequest(http.MethodGet, "http://api.marketstack.com/v1/tickers/"+symbol+"/eod?access_key=214e04d7d20abb7cbd47894b00c2c334")
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
	log.Print(val)
	log.Print("body", string(req))
	return &val, nil
}
