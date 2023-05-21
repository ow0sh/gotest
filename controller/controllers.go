package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/config"
)

func GetCoins(client *http.Client) map[string]string {
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/list", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var tmpCoinlist []config.Coin
	json.NewDecoder(resp.Body).Decode(&tmpCoinlist)

	result := make(map[string]string)
	for _, v := range tmpCoinlist {
		result[v.Symbol] = v.Name
	}

	return result
}

func GetSupported(client *http.Client) []string {
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/supported_vs_currencies", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var tmpSupportedCoinlist []string
	json.NewDecoder(resp.Body).Decode(&tmpSupportedCoinlist)

	return tmpSupportedCoinlist
}

func GetPrice(base, quote string) json.RawMessage {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%v&vs_currencies=%v", base, quote), nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var jsonPrice json.RawMessage
	json.NewDecoder(resp.Body).Decode(&jsonPrice)

	return jsonPrice
}

func Contains(s interface{}, str string) bool {
	switch value := s.(type) {
	case []string:
		for _, v := range value {
			if v == str {
				return true
			}
		}

		return false
	case map[string]string:
		for k := range value {
			if k == str {
				return true
			}
		}

		return false
	}
	return false
}
