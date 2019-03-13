package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Request is a collection containing an endpoint, method, and headers. It's essentially an abstraction of an http request
type Request struct {
	Endpoint string
	Method   string
	Body     string
	Headers  map[string]string
}

// SetHeader sets the specified header
func (r *Request) SetHeader(key, val string) {
	r.Headers[key] = val
}

func (r *Request) Get() (string, error) {
	r.Method = "GET"
	return r.Do()
}

func (r *Request) Post(vals map[string]string) (string, error) {
	r.Method = "POST"

	body, _ := json.Marshal(vals)
	r.Body = string(body)

	return r.Do()
}

func (r *Request) Do() (string, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(r.Method, r.Endpoint, strings.NewReader(r.Body))
	if err != nil {
		return "", errors.Wrap(err, "http request initialization failed")
	}

	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "http request execution failed")
	}

	respString, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", errors.Wrap(err, "http response reading failed")
	}

	return string(respString), nil
}
