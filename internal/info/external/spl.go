package external

import (
	"fmt"
	"net/url"

	"github.com/trustwallet/assets-go-libs/http"
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
	apiURL := fmt.Sprintf(splAPIURL, url.QueryEscape(tokenID))

	var result TokenInfoSPL
	err := http.GetHTTPResponse(apiURL, &result)
	if err != nil {
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
