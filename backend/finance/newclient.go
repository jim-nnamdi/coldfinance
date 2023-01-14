package finance

import (
	"net/http"
	"net/url"
)

func Mrequest(r *http.Request, path string) (*http.Request, error) {
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

func Dorequest(req *http.Request) (*http.Response, error) {
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return res, err
}
