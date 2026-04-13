package external

import (
	"fmt"
)

const splAPIURL = "https://public-api.solscan.io/token/holders?tokenAddress=%s"

type TokenInfoSPL struct {
	Data         []Data `json:"data"`
	HoldersCount int    `json:"total"`
}

type Data struct {
	Decimals int `json:"decimals"`
}

func GetTokenInfoForSPL(tokenID string) (*TokenInfo, error) {
	url := fmt.Sprintf(splAPIURL, tokenID)

	var result TokenInfoSPL
	if err := getJSON(url, &result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("failed to get token info for SPL token")
	}

	return &TokenInfo{
		Decimals:     result.Data[0].Decimals,
		HoldersCount: result.HoldersCount,
	}, nil
}
