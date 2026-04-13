package external

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// ethplorerAPIKeyEnv is the environment variable used to override the default
// Ethplorer API key. When unset, the rate-limited public "freekey" is used.
const ethplorerAPIKeyEnv = "ETHPLORER_API_KEY"

func ethplorerAPIKey() string {
	if key := os.Getenv(ethplorerAPIKeyEnv); key != "" {
		return key
	}

	return "freekey"
}

type TokenInfoERC20 struct {
	Decimals     string `json:"decimals"`
	HoldersCount int    `json:"holdersCount"`
}

func GetTokenInfoForERC20(tokenID string) (*TokenInfo, error) {
	apiURL := fmt.Sprintf("https://api.ethplorer.io/getTokenInfo/%s?apiKey=%s",
		url.PathEscape(tokenID), url.QueryEscape(ethplorerAPIKey()))

	var result TokenInfoERC20
	if err := getJSON(apiURL, &result); err != nil {
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
