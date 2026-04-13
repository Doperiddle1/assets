package external

import (
	"fmt"
	"net/url"
 copilot/fix-security-issues


	"github.com/trustwallet/assets-go-libs/http"
 master
)

type TokenInfoSPL struct {
	Data         []Data `json:"data"`
	HoldersCount int    `json:"total"`
}

type Data struct {
	Decimals int `json:"decimals"`
}

func GetTokenInfoForSPL(tokenID string) (*TokenInfo, error) {
 copilot/fix-security-issues
	apiURL := "https://public-api.solscan.io/token/holders?tokenAddress=" + url.QueryEscape(tokenID)

	var result TokenInfoSPL
	if err := getJSON(apiURL, &result); err != nil {

	apiURL := fmt.Sprintf(splAPIURL, url.QueryEscape(tokenID))

	var result TokenInfoSPL
	err := http.GetHTTPResponse(apiURL, &result)
	if err != nil {
 master
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
