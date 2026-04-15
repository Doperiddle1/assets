package external

import (
	"errors"
	"net/url"
)

type TRC20TokensResponse struct {
	TRC20Tokens []struct {
		Symbol       string `json:"symbol"`
		Decimals     int    `json:"decimals"`
		HoldersCount int    `json:"holders_count"`
	} `json:"trc20_tokens"`
}

func GetTokenInfoForTRC20(tokenID string) (*TokenInfo, error) {
	apiURL := "https://apilist.tronscan.io/api/token_trc20?contract=" + url.QueryEscape(tokenID)

	var res TRC20TokensResponse
	if err := getJSON(apiURL, &res); err != nil {
		return nil, err
	}

	if len(res.TRC20Tokens) == 0 {
		return nil, errors.New("not found")
	}

	return &TokenInfo{
		Symbol:       res.TRC20Tokens[0].Symbol,
		Decimals:     res.TRC20Tokens[0].Decimals,
		HoldersCount: res.TRC20Tokens[0].HoldersCount,
	}, nil
}
