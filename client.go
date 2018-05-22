package example

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	baseURL = "https://some.com/api/v1"
)

type Option func(*Client)

func SetHTTPClient(httpClient *http.Client) Option {
	return func(cli *Client) {
		cli.httpClient = httpClient
	}
}

type Client struct {
	api, secret string
	httpClient  *http.Client
}

func NewClient(api, secret string, options ...Option) *Client {
	cli := Client{
		api:    api,
		secret: secret,
		httpClient: &http.Client{
			Timeout: time.Second,
		},
	}

	for i := range options {
		options[i](&cli)
	}

	return &cli
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (cli *Client) GetUsers() ([]User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users", baseURL), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build request")
	}

	resp, err := cli.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	res := struct {
		Users []User `json:""`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, errors.Wrap(err, "unmarshaling failed")
	}

	return res.Users, nil

}
