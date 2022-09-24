package main

import (
	"testing"
)

const (
	INPUT = `-----BEGIN CERTIFICATE-----
some-contents
more contents --
-----END CERTIFICATE-----`

	EXPECTED_OUTPUT = `some-contents
more contents --`
)

func TestRemoveHeaderTrailer(t *testing.T) {
	t.Run("Remove header and trailer", func(t *testing.T) {
		result := removeHeaderTrailer([]byte(INPUT))
		if result != EXPECTED_OUTPUT {
			t.Errorf("got %s, want %s", result, EXPECTED_OUTPUT)
		}
	})
}

func TestVerifyUnsignedToken(t *testing.T) {
	var input = "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJBdHRlc3RhdGlvblBvbGljeSI6IlpHVm1ZWFZzZENCdFlYUmphR1Z6WDNObmVGOXdiMnhwWTNrZ1BTQm1ZV3h6WlFwdFlYUmphR1Z6WDNObmVGOXdiMnhwWTNrZ1BTQjBjblZsSUhzS0lDQWdhVzV3ZFhRdVlXMWlaWEpmYzJkNFgyMXlaVzVqYkdGMlpTQTlQU0FpT0RObU5HVTRNVGs0TmpGaFpHVm1ObVptWWpKaE5EZzJOV1ZtWldFNU16TTNZamt4WldRek1HWmhNek0wT1RGaU1UZG1NR1ExWkRsbE9ESXdORFF4TUNJS0lDQWdhVzV3ZFhRdVlXMWlaWEpmYzJkNFgyMXljMmxuYm1WeUlEMDlJQ0k0TTJRM01UbGxOemRrWldGallURTBOekJtTm1KaFpqWXlZVFJrTnpjME16QXpZemc1T1dSaU5qa3dNakJtT1dNM01HVmxNV1JtWXpBNFl6ZGpaVGxsSWdvZ0lDQnBibkIxZEM1aGJXSmxjbDl6WjNoZmFYTjJjSEp2Wkdsa0lEMDlJREFLSUNBZ2FXNXdkWFF1WVcxaVpYSmZjMmQ0WDJsemRuTjJiaUE5UFNBd0NpQWdJR2x1Y0hWMExtRnRZbVZ5WDNObmVGOXBjMTlrWldKMVoyZGhZbXhsSUQwOUlHWmhiSE5sQ24wSyJ9."
	t.Run("Verify unsigned policy token succeed", func(t *testing.T) {
		verified, err := verifyToken(input)
		if err != nil {
			t.Errorf("got err %s", err.Error())
		}
		if verified != true {
			t.Errorf("got verified %v expect true", verified)
		}
	})

	t.Run("Verify unsigned policy token fail", func(t *testing.T) {
		verified, err := verifyToken(input + "a")
		if err != nil {
			t.Logf("got err %s", err.Error())
		} else {
			t.Error("Expect error to happen")
		}
		if verified == true {
			t.Errorf("got verified %v expect false", verified)
		}
	})
}
