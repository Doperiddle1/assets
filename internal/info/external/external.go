package external

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	holdersRegexp  = regexp.MustCompile(`(\d+)\saddresses`)
	decimalsRegexp = regexp.MustCompile(`(\d+)\s+<\/div>`)
	symbolRegexp   = regexp.MustCompile(`<b>(\w+)<\/b>\s<span`)

	// evmAddressRegexp matches a standard 0x-prefixed EVM address (20 bytes, hex).
	evmAddressRegexp = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	// trc20AddressRegexp matches a TRON address (base58, starts with T, 34 chars).
	trc20AddressRegexp = regexp.MustCompile(`^T[1-9A-HJ-NP-Za-km-z]{33}$`)
	// trc10IDRegexp matches a TRON TRC10 numeric token ID.
	trc10IDRegexp = regexp.MustCompile(`^\d+$`)
	// splAddressRegexp matches a Solana base58 public key (32 bytes → 43–44 base58 chars).
	splAddressRegexp = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{43,44}$`)
)

type TokenInfo struct {
	Symbol       string
	Decimals     int
	HoldersCount int
}

func GetTokenInfo(tokenID, tokentType string) (*TokenInfo, error) {
	escapedID := url.PathEscape(tokenID)

	switch strings.ToLower(tokentType) {
	case "erc20":
		if !evmAddressRegexp.MatchString(tokenID) {
			return nil, fmt.Errorf("invalid ERC20 token address: %q", tokenID)
		}
		return GetTokenInfoForERC20(tokenID)
	case "bep20":
		if !evmAddressRegexp.MatchString(tokenID) {
			return nil, fmt.Errorf("invalid BEP20 token address: %q", tokenID)
		}
		return GetTokenInfoByScraping(fmt.Sprintf("https://bscscan.com/token/%s", escapedID))
	case "fantom":
		if !evmAddressRegexp.MatchString(tokenID) {
			return nil, fmt.Errorf("invalid Fantom token address: %q", tokenID)
		}
		return GetTokenInfoByScraping(fmt.Sprintf("https://ftmscan.com/token/%s", escapedID))
	case "polygon":
		if !evmAddressRegexp.MatchString(tokenID) {
			return nil, fmt.Errorf("invalid Polygon token address: %q", tokenID)
		}
		return GetTokenInfoByScraping(fmt.Sprintf("https://polygonscan.com/token/%s", escapedID))
	case "avalanche":
		if !evmAddressRegexp.MatchString(tokenID) {
			return nil, fmt.Errorf("invalid Avalanche token address: %q", tokenID)
		}
		return GetTokenInfoByScraping(fmt.Sprintf("https://snowtrace.io/token/%s", escapedID))
	case "spl":
		if !splAddressRegexp.MatchString(tokenID) {
			return nil, fmt.Errorf("invalid SPL token address: %q", tokenID)
		}
		return GetTokenInfoForSPL(tokenID)
	case "trc20":
		if !trc20AddressRegexp.MatchString(tokenID) {
			return nil, fmt.Errorf("invalid TRC20 token address: %q", tokenID)
		}
		return GetTokenInfoForTRC20(tokenID)
	case "trc10":
		if !trc10IDRegexp.MatchString(tokenID) {
			return nil, fmt.Errorf("invalid TRC10 token ID: %q", tokenID)
		}
		return GetTokenInfoForTRC10(tokenID)
	}

	return nil, nil
}

func GetTokenInfoByScraping(url string) (*TokenInfo, error) {
	data, err := getBytes(url)
	if err != nil {
		return nil, err
	}

	// Remove all "," from content.
	pageContent := strings.ReplaceAll(string(data), ",", "")

	var holders, decimals int
	var symbol string

	match := symbolRegexp.FindStringSubmatch(pageContent)
	if len(match) > 1 {
		symbol = match[1]
	}

	match = decimalsRegexp.FindStringSubmatch(pageContent)
	if len(match) > 1 {
		decimals, err = strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}
	}

	match = holdersRegexp.FindStringSubmatch(pageContent)
	if len(match) > 1 {
		holders, err = strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}
	}

	return &TokenInfo{
		Symbol:       symbol,
		Decimals:     decimals,
		HoldersCount: holders,
	}, nil
}
