package external

import (
	"fmt"
	"net/url"
 copilot/fix-security-issues-plan

 copilot/fix-security-issues
	"os"

 master
 master
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

// maxTokenDecimals is the maximum number of decimal places a valid token can have.
const maxTokenDecimals = 30

type TokenInfoERC20 struct {
	Decimals     string `json:"decimals"`
	HoldersCount int    `json:"holdersCount"`
}

func GetTokenInfoForERC20(tokenID string) (*TokenInfo, error) {
 copilot/fix-security-issues-plan
	url := fmt.Sprintf(ethAPIURL, url.PathEscape(tokenID))

	var result TokenInfoERC20
	err := getHTTPResponse(url, &result)

 copilot/fix-security-issues
	apiURL := fmt.Sprintf("https://api.ethplorer.io/getTokenInfo/%s?apiKey=%s",
		url.PathEscape(tokenID), url.QueryEscape(ethplorerAPIKey()))

	var result TokenInfoERC20
	if err := getJSON(apiURL, &result); err != nil {

	apiURL := fmt.Sprintf(ethAPIURL, url.PathEscape(tokenID))

	var result TokenInfoERC20
	err := http.GetHTTPResponse(apiURL, &result)
 master
	if err != nil {
 master
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
