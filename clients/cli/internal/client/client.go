package client

import (
	"net/http"
)


const BASEURL = "http://localhost:9847"


type Client struct {
	HTTP    *http.Client
}

func New() *Client {
	return &Client{
		HTTP:    &http.Client{},
	}
}


func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", BASEURL + url, nil)
	if err != nil {
		return nil, err
	}

	return c.HTTP.Do(req)
}
