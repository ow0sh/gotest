package coingecko

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ow0sh/gotest/pkg/httpconn"
	"github.com/pkg/errors"
)

const BASE_URL = "https://api.coingecko.com"

type Client struct {
	cn *httpconn.Connector
}

func NewClient() *Client {
	u, _ := url.Parse(BASE_URL)
	timeout := time.Second * 10
	return &Client{
		cn: httpconn.NewConnector(u, &timeout),
	}
}

func (cli *Client) GetCoinList(ctx context.Context) (map[string]string, error) {
	path := "/api/v3/coins/list"
	statusCode, body, err := cli.cn.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	if *statusCode == http.StatusOK {
		var resp []Coin
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal responce")
		}
		result := make(map[string]string)
		for _, v := range resp {
			result[v.Symbol] = v.Name
		}
		return result, nil
	}

	return nil, errors.New(fmt.Sprintf("failed with  %s %s", http.StatusText(*statusCode), body))
}

func (cli *Client) GetSupported(ctx context.Context) ([]string, error) {
	path := "/api/v3/simple/supported_vs_currencies"
	statusCode, body, err := cli.cn.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	if *statusCode == http.StatusOK {
		var resp []string
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal responce")
		}
		return resp, nil
	}

	return nil, errors.New(fmt.Sprintf("failed with  %s %s", http.StatusText(*statusCode), body))
}

func (cli *Client) GetPrice(ctx context.Context, base string, quote string) (json.RawMessage, error) {
	path := fmt.Sprintf("/api/v3/simple/price?ids=%v&vs_currencies=%v", base, quote)
	statusCode, body, err := cli.cn.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	if *statusCode == http.StatusOK {
		var jsonPrice json.RawMessage
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&jsonPrice); err != nil {
			return nil, errors.Wrap(err, "failed to decode json")
		}
		return jsonPrice, nil
	}

	return nil, errors.New(fmt.Sprintf("failed with  %s %s", http.StatusText(*statusCode), body))
}
