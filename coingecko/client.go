package coingecko

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type Client struct {
	httpCli *http.Client
}

type Coin struct {
	Id     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func NewClient(cli *http.Client) *Client {
	return &Client{httpCli: cli}
}

func (cli *Client) GetCoins() (map[string]string, error) {
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/list", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := cli.httpCli.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}
	defer resp.Body.Close()

	var tmpCoinlist []Coin
	if err = json.NewDecoder(resp.Body).Decode(&tmpCoinlist); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	result := make(map[string]string)
	for _, v := range tmpCoinlist {
		result[v.Symbol] = v.Name
	}

	return result, nil
}

func (cli *Client) GetSupported() []string {
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/supported_vs_currencies", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := cli.httpCli.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var tmpSupportedCoinlist []string
	json.NewDecoder(resp.Body).Decode(&tmpSupportedCoinlist)

	return tmpSupportedCoinlist
}

func (cli *Client) GetPrice(base, quote string) json.RawMessage {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%v&vs_currencies=%v", base, quote), nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := cli.httpCli.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var jsonPrice json.RawMessage
	json.NewDecoder(resp.Body).Decode(&jsonPrice)

	return jsonPrice
}
