package processor

import (
	"testing"

	"github.com/trustwallet/assets-go-libs/validation"
)

func TestCalculateTargetDimension(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		width  int
		height int
		wantW  int
		wantH  int
	}{
		{
			name:   "square already at max",
			width:  512,
			height: 512,
			wantW:  512,
			wantH:  512,
		},
		{
			name:   "square larger than max",
			width:  1024,
			height: 1024,
			wantW:  512,
			wantH:  512,
		},
		{
			name:   "landscape — width is the max edge",
			width:  1024,
			height: 512,
			wantW:  512,
			wantH:  256,
		},
		{
			name:   "portrait — height is the max edge",
			width:  256,
			height: 1024,
			wantW:  128,
			wantH:  512,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotW, gotH := calculateTargetDimension(tc.width, tc.height)
			if gotW != tc.wantW || gotH != tc.wantH {
				t.Errorf("calculateTargetDimension(%d, %d) = (%d, %d), want (%d, %d)",
					tc.width, tc.height, gotW, gotH, tc.wantW, tc.wantH)
			}
			// Ensure neither dimension exceeds the allowed maximum.
			if gotW > validation.MaxW || gotH > validation.MaxH {
				t.Errorf("result (%d x %d) exceeds max (%d x %d)", gotW, gotH, validation.MaxW, validation.MaxH)
			}
		})
	}
}
