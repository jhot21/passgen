package passgen

import (
	"strings"
	"testing"
)

func TestGuards(t *testing.T) {
	cases := []struct {
		name string
		b    *Builder
	}{
		{"v2 alphabet of one char", New("p").WithVersion(2).NoUppers().NoNumbers().NoSpecials().
			WithAmbiguousChars("bcdefghijklmnopqrstuvwxyz")},
		{"v2 alphabet empty", New("p").WithVersion(2).NoUppers().NoNumbers().NoSpecials().
			WithAmbiguousChars("abcdefghijklmnopqrstuvwxyz")},
		{"v2 specials filtered away", New("p").WithVersion(2).WithSpecials("l").WithAmbiguousChars("l")},
		{"v2 uppers filtered away", New("p").WithVersion(2).WithAmbiguousChars("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
		{"v2 length shorter than required classes", New("p").WithVersion(2).WithLength(2)},
		{"v3 class longer than 256", New("p").WithVersion(3).WithSpecials(strings.Repeat("!", 257))},
		{"v3 all classes filtered away", New("p").WithVersion(3).NoUppers().NoNumbers().NoSpecials().
			WithAmbiguousChars("abcdefghijklmnopqrstuvwxyz")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.b.Generate(); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}

func TestGuardsDoNotRejectValidInputs(t *testing.T) {
	valid := []*Builder{
		New("p").WithVersion(2).WithLength(3),
		New("p").WithVersion(2).NoAmbiguous(),
		New("p").WithVersion(3).WithSpecials(strings.Repeat("!", 256)),
		New("p").WithVersion(3).WithLength(1),
	}
	for i, b := range valid {
		if _, err := b.Generate(); err != nil {
			t.Fatalf("case %d: unexpected error: %v", i, err)
		}
	}
}
