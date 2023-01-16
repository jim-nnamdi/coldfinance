package web

import (
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

type NewClient interface {
	Mrequest(r *http.Request, path string) (*http.Request, error)
	Dorequest(req *http.Request) (*http.Response, error)
}

var _ NewClient = &NewC{}

type NewC struct {
	logger     *zap.Logger
	httpClient *http.Client
}

func NewClients(logger *zap.Logger, httpclient *http.Client) *NewC {
	return &NewC{
		logger:     logger,
		httpClient: httpclient,
	}
}

func (nc *NewC) Mrequest(r *http.Request, path string) (*http.Request, error) {
	var (
		newreq *http.Request
		err    error
	)
	if newreq, err = http.NewRequestWithContext(r.Context(), r.Method, path, r.Body); err != nil {
		return nil, err
	}
	for k, vv := range r.Header {
		for _, v := range vv {
			newreq.Header.Set(k, v)
		}
	}
	newreq.PostForm = url.Values{}
	for k, vv := range r.Form {
		for _, v := range vv {
			newreq.Form.Set(k, v)
		}
	}
	newreq.Header.Add("Content-Type", "application/json")
	return newreq, nil
}

func (nc *NewC) Dorequest(req *http.Request) (*http.Response, error) {
	res, err := nc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, err
}
