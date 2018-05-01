package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Request is a collection containing an endpoint, method, and headers. It's essentially an abstraction of an http request
type Request struct {
	endpoint string
	method   string
	headers  map[string]string
}

// SetHeader sets the specified header
func (r *Request) SetHeader(key, val string) {
	r.headers[key] = val
}

// Do executes the http request & returns a PageCollection and an error
func (r *Request) Do() (*PageCollection, error) {
	pages := &PageCollection{}
	hasNext := true
	endpoint := r.endpoint

	for hasNext {
		log.Printf("[API] %s %s\n", r.method, endpoint)

		httpClient := &http.Client{
			Timeout: time.Second * 10,
		}

		req, err := http.NewRequest(r.method, endpoint, nil)
		if err != nil {
			return nil, err
		}

		for k, v := range r.headers {
			req.Header.Add(k, v)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		respString, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}

		page, err := ParsePage(string(respString))
		if err != nil {
			return nil, err
		}

		pages.Add(page)

		endpoint, hasNext = page["next"].(string)
	}

	return pages, nil
}
