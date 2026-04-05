package external

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/trustwallet/assets-go-libs/http"
)

const ethAPIURL = "https://api.ethplorer.io/getTokenInfo/%s?apiKey=freekey"

// maxTokenDecimals is the maximum number of decimal places a valid token can have.
const maxTokenDecimals = 30

type TokenInfoERC20 struct {
	Decimals     string `json:"decimals"`
	HoldersCount int    `json:"holdersCount"`
}

func GetTokenInfoForERC20(tokenID string) (*TokenInfo, error) {
	apiURL := fmt.Sprintf(ethAPIURL, url.PathEscape(tokenID))

	var result TokenInfoERC20
	err := http.GetHTTPResponse(apiURL, &result)
	if err != nil {
		return nil, err
	}

	decimals, err := strconv.Atoi(result.Decimals)
	if err != nil {
		return nil, err
	}

	if decimals < 0 || decimals > maxTokenDecimals {
		return nil, fmt.Errorf("decimals value out of valid range: %d", decimals)
	}

	return &TokenInfo{
		Decimals:     decimals,
		HoldersCount: result.HoldersCount,
	}, nil
}
