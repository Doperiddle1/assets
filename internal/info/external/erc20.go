package external

import (
	"fmt"
	"strconv"
)

const (
	ethplorerAPIKey = "freekey"
	ethAPIURL       = "https://api.ethplorer.io/getTokenInfo/%s?apiKey=" + ethplorerAPIKey
)

type TokenInfoERC20 struct {
	Decimals     string `json:"decimals"`
	HoldersCount int    `json:"holdersCount"`
}

func GetTokenInfoForERC20(tokenID string) (*TokenInfo, error) {
	url := fmt.Sprintf(ethAPIURL, tokenID)

	var result TokenInfoERC20
	err := getJSON(url, &result)
	if err != nil {
		return nil, err
	}

	decimals, err := strconv.Atoi(result.Decimals)
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		Decimals:     decimals,
		HoldersCount: result.HoldersCount,
	}, nil
}
