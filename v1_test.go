package passgen

import (
	"testing"
	"unicode/utf16"
)

func TestPassgenV1Regression(t *testing.T) {
	got := passgenV1(opts{passphrase: "mypassphrase", salt: "testsalt", length: 40})
	want := "*GS3t`E,*9XR29/Wm7jhw&1:mJB}RI;8YbR{RL.+"
	if got != want {
		t.Errorf("passgenV1 = %q, want %q", got, want)
	}
}

func TestPassgenV1EmptySalt(t *testing.T) {
	got := passgenV1(opts{passphrase: "mypassphrase", salt: "", length: 40})
	want := "1t~IAk<0<1V<9*AS9Oj.xl?Hp!&u_jGS)GX<ll9X"
	if got != want {
		t.Errorf("passgenV1 empty salt = %q, want %q", got, want)
	}
}

func TestConvertUTF16ToLittleEndianBytes(t *testing.T) {
	// "hi" in UTF-16: h=0x0068, i=0x0069
	input := utf16.Encode([]rune("hi"))
	got := convertUTF16ToLittleEndianBytes(input)
	want := []byte{0x68, 0x00, 0x69, 0x00}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("byte[%d] = 0x%02x, want 0x%02x", i, got[i], want[i])
		}
	}
}

func TestHexToInt(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"0", 0},
		{"9", 9},
		{"a", 10},
		{"f", 15},
		{"F", 15},
		{"A", 10},
	}
	for _, tt := range tests {
		got := hexToInt(tt.input)
		if got != tt.want {
			t.Errorf("hexToInt(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
