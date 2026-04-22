package info_test

import (
	"strings"
	"testing"

	"github.com/trustwallet/assets/internal/info"
)

// ---------- ValidateDecimals ----------

func TestValidateDecimals(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		decimals int
		wantErr  bool
	}{
		{"zero", 0, false},
		{"typical ERC20", 18, false},
		{"max allowed", 30, false},
		{"negative", -1, true},
		{"over max", 31, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := info.ValidateDecimals(tc.decimals)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateDecimals(%d) error = %v, wantErr %v", tc.decimals, err, tc.wantErr)
			}
		})
	}
}

// ---------- ValidateStatus ----------

func TestValidateStatus(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{"active", "active", false},
		{"spam", "spam", false},
		{"abandoned", "abandoned", false},
		{"empty", "", true},
		{"unknown", "pending", true},
		{"mixed case", "Active", true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := info.ValidateStatus(tc.status)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateStatus(%q) error = %v, wantErr %v", tc.status, err, tc.wantErr)
			}
		})
	}
}

// ---------- ValidateDescription ----------

func TestValidateDescription(t *testing.T) {
	t.Parallel()

	longDesc := strings.Repeat("a", 601)

	cases := []struct {
		name    string
		desc    string
		wantErr bool
	}{
		{"empty", "", false},
		{"normal", "A simple token.", false},
		{"max length", strings.Repeat("a", 600), false},
		{"over max length", longDesc, true},
		{"contains newline", "line1\nline2", true},
		{"contains double space", "too  many spaces", true},
		{"contains HTML tag", "hello <b>world</b>", true},
		{"self-closing tag", "<br/>", true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := info.ValidateDescription(tc.desc)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateDescription(%q) error = %v, wantErr %v", tc.desc, err, tc.wantErr)
			}
		})
	}
}

// ---------- ValidateDescriptionWebsite ----------

func TestValidateDescriptionWebsite(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		desc    string
		website string
		wantErr bool
	}{
		{"description with website", "A token.", "https://example.com", false},
		{"dash description no website", "-", "", false},
		{"description without website", "A token.", "", true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := info.ValidateDescriptionWebsite(tc.desc, tc.website)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateDescriptionWebsite(%q, %q) error = %v, wantErr %v",
					tc.desc, tc.website, err, tc.wantErr)
			}
		})
	}
}

// ---------- ValidateAssetID ----------

func TestValidateAssetID(t *testing.T) {
	t.Parallel()

	addr := "0xAbCdEf1234567890abcdef1234567890AbCdEf12"

	cases := []struct {
		name    string
		id      string
		address string
		wantErr bool
	}{
		{"exact match", addr, addr, false},
		{"different string", "0xOther", addr, true},
		{"case mismatch only", strings.ToLower(addr), addr, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := info.ValidateAssetID(tc.id, tc.address)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateAssetID(%q, %q) error = %v, wantErr %v",
					tc.id, tc.address, err, tc.wantErr)
			}
		})
	}
}

// ---------- ValidateLinks ----------

func TestValidateLinks(t *testing.T) {
	t.Parallel()

	strPtr := func(s string) *string { return &s }

	cases := []struct {
		name    string
		links   []info.Link
		wantErr bool
	}{
		{"nil links", nil, false},
		{"single link (< 2, skip check)", []info.Link{{Name: strPtr("github"), URL: strPtr("https://github.com/foo")}}, false},
		{"two valid links", []info.Link{
			{Name: strPtr("github"), URL: strPtr("https://github.com/foo")},
			{Name: strPtr("twitter"), URL: strPtr("https://twitter.com/foo")},
		}, false},
		{"unknown link name", []info.Link{
			{Name: strPtr("github"), URL: strPtr("https://github.com/foo")},
			{Name: strPtr("unknown_link"), URL: strPtr("https://example.com/foo")},
		}, true},
		{"http instead of https", []info.Link{
			{Name: strPtr("github"), URL: strPtr("https://github.com/foo")},
			{Name: strPtr("whitepaper"), URL: strPtr("http://example.com/paper.pdf")},
		}, true},
		{"github wrong prefix", []info.Link{
			{Name: strPtr("github"), URL: strPtr("https://gitlab.com/foo")},
			{Name: strPtr("twitter"), URL: strPtr("https://twitter.com/foo")},
		}, true},
		{"medium without medium.com", []info.Link{
			{Name: strPtr("github"), URL: strPtr("https://github.com/foo")},
			{Name: strPtr("medium"), URL: strPtr("https://blog.example.org/article")},
		}, true},
		{"nil name", []info.Link{
			{Name: nil, URL: strPtr("https://github.com/foo")},
			{Name: strPtr("twitter"), URL: strPtr("https://twitter.com/foo")},
		}, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := info.ValidateLinks(tc.links)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateLinks error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// ---------- ValidateCoinType ----------

func TestValidateCoinType(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		assetType string
		wantErr   bool
	}{
		{"valid", "coin", false},
		{"token", "token", true},
		{"mixed case", "Coin", true},
		{"empty", "", true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := info.ValidateCoinType(tc.assetType)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateCoinType(%q) error = %v, wantErr %v", tc.assetType, err, tc.wantErr)
			}
		})
	}
}

// ---------- ValidateTags ----------

func TestValidateTags(t *testing.T) {
	t.Parallel()

	allowed := []string{"defi", "stablecoin", "wrapped"}

	cases := []struct {
		name    string
		tags    []string
		wantErr bool
	}{
		{"no tags", nil, false},
		{"all allowed", []string{"defi", "stablecoin"}, false},
		{"one unknown", []string{"defi", "nft"}, true},
		{"all unknown", []string{"foo", "bar"}, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := info.ValidateTags(tc.tags, allowed)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateTags(%v) error = %v, wantErr %v", tc.tags, err, tc.wantErr)
			}
		})
	}
}

// ---------- ValidateAssetDecimalsAccordingType ----------

func TestValidateAssetDecimalsAccordingType(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		assetType string
		decimals  int
		wantErr   bool
	}{
		{"BEP2 correct", "BEP2", 8, false},
		{"BEP2 wrong", "BEP2", 18, true},
		{"ERC20 any decimals", "ERC20", 18, false},
		{"ERC20 zero decimals", "ERC20", 0, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := info.ValidateAssetDecimalsAccordingType(tc.assetType, tc.decimals)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateAssetDecimalsAccordingType(%q, %d) error = %v, wantErr %v",
					tc.assetType, tc.decimals, err, tc.wantErr)
			}
		})
	}
}
