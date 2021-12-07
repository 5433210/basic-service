package http

import "net/http"

type StdHttpRequestDoer struct{}

func (d StdHttpRequestDoer) Do(req *http.Request) (*http.Response, error) {
	return nil, nil
}
