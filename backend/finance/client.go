package finance

import (
	"io"
	"net/http"

	"go.uber.org/zap"
)

type RequestClient interface {
	MakeGetRequest(method string, url string) ([]byte, error)
}

var _ RequestClient = &Dataclient{}

type Dataclient struct {
	logger *zap.Logger
}

func NewDataClient(logger *zap.Logger) *Dataclient {
	return &Dataclient{
		logger: logger,
	}
}

func (c *Dataclient) MakeGetRequest(method string, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		c.logger.Debug("error making get request", zap.Any("error", err.Error()))
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.logger.Debug("error processing data", zap.Any("error", err.Error()))
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return body, nil
}
