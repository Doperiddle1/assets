package external

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/trustwallet/assets-go-libs/http"
)

const trc10APIURL = "https://apilist.tronscan.io/api/token?id=%s"

type TRC10TokensResponse struct {
	Data []struct {
		Symbol       string `json:"abbr"`
		Decimals     int    `json:"precision"`
		HoldersCount int    `json:"nrOfTokenHolders"`
	} `json:"data"`
}

func GetTokenInfoForTRC10(tokenID string) (*TokenInfo, error) {
	apiURL := fmt.Sprintf(trc10APIURL, url.QueryEscape(tokenID))

	var res TRC10TokensResponse
	err := http.GetHTTPResponse(apiURL, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Data) == 0 {
		return nil, errors.New("not found")
	}

	return &TokenInfo{
		Symbol:       res.Data[0].Symbol,
		Decimals:     res.Data[0].Decimals,
		HoldersCount: res.Data[0].HoldersCount,
	}, nil
}
