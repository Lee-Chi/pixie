package http

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewRequest() *Request {
	return &Request{
		header: map[string]string{},
	}
}

type Request struct {
	header map[string]string
}

func (r *Request) SetHeader(key string, value string) *Request {
	r.header[key] = value
	return r
}

func (r Request) Get(api string) ([]byte, error) {
	client := http.Client{
		Timeout: 300 * time.Second,
	}

	request, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range r.header {
		request.Header.Set(k, v)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is %d", response.StatusCode)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
