package finance

import (
	"encoding/json"
	"log"
	"net/http"

	"go.uber.org/zap"
)

type CryptoData interface {
	GetAllCryptoData() (*AllCrypto, error)
	GetLiveCryptoData() (*LiveData, error)
	ConvertCrypto(coinfrom string, cointo string, amount int) (float64, error)
}

var _ CryptoData = &Cryptos{}

type Cryptos struct {
	logger  *zap.Logger
	sclient RequestClient
}

type CryptoResponse struct {
	Success bool `json:"success"`
	Crypto  struct {
	} `json:"crypto"`
}

type AllCrypto struct {
	Success bool        `json:"success"`
	Crypto  interface{} `json:"crypto"`
}

type LiveData struct {
	Timestamp int `json:"timestamp"`

	// currency we want to display
	// the current rates in
	Target string             `json:"target"`
	Rates  map[string]float64 `json:"rates"`
}

func NewCryptos(logger *zap.Logger, sclient RequestClient) *Cryptos {
	return &Cryptos{
		logger:  logger,
		sclient: sclient,
	}
}

func (cs *Cryptos) GetAllCryptoData() (*AllCrypto, error) {
	req, err := cs.sclient.MakeGetRequest(http.MethodGet, "http://api.coinlayer.com/api/list?access_key=09e5664c792fbfa119e714e95095eedc")
	if err != nil {
		cs.logger.Debug("error fetching data", zap.Any("error", err))
		return nil, err
	}
	var val AllCrypto
	_ = json.Unmarshal(req, &val)
	return &val, nil
}

func (cs *Cryptos) GetLiveCryptoData() (*LiveData, error) {
	req, err := cs.sclient.MakeGetRequest(http.MethodGet, "http://api.coinlayer.com/api/live?access_key=09e5664c792fbfa119e714e95095eedc")
	if err != nil {
		cs.logger.Debug("error fetching data", zap.Any("error", err))
		return nil, err
	}
	var val LiveData
	_ = json.Unmarshal(req, &val)
	return &val, nil
}

func (cs *Cryptos) ConvertCrypto(coinfrom string, cointo string, amount int) (float64, error) {
	req, err := cs.sclient.MakeGetRequest(http.MethodGet, "http://api.coinlayer.com/api/live?access_key=09e5664c792fbfa119e714e95095eedc")
	if err != nil {
		cs.logger.Debug("error fetching data", zap.Any("error", err))
		return 0.0, err
	}
	var val LiveData
	_ = json.Unmarshal(req, &val)
	if val.Rates != nil {
		var coinFromPrice float64
		var coinToPrice float64
		for k, v := range val.Rates {
			if k == coinfrom {
				log.Print("coinfrom price -> ", v)
				coinFromPrice = v
			}
			if k == cointo {
				log.Print("cointo price -> ", v)
				coinToPrice = v
			}
		}
		calcCoinFromAndAmount := coinFromPrice * float64(amount)
		calcConversionPrice := calcCoinFromAndAmount / coinToPrice
		return calcConversionPrice, nil
	}
	return 0.0, nil
}
