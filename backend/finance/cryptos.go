package finance

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type CryptoData interface {
	GetAllCryptoData()
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

func NewCryptos(logger *zap.Logger, sclient RequestClient) *Cryptos {
	return &Cryptos{
		logger:  logger,
		sclient: sclient,
	}
}

func (cx *Cryptos) GetAllCryptoData() {
	req, err := cx.sclient.MakeGetRequest(http.MethodGet, "http://api.coinlayer.com/api/list?access_key=09e5664c792fbfa119e714e95095eedc")
	if err != nil {
		cx.logger.Debug("error fetching data", zap.Any("error", err))
		return
	}
	fmt.Print(req)
}
