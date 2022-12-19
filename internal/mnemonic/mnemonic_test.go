package mnemonic_test

import (
	"fmt"
	"testing"

	"mnemonic.deflinhec.dev/internal/mnemonic"
)

func Test(t *testing.T) {
	fetcher := mnemonic.Fetcher().Fetch(`TXaMXTQgtdV6iqxtmQ7HNnqzXRoJKXfFAz`,
		mnemonic.Pharse(
			`range sheriff try enroll deer over ten level bring display stamp`,
		)).Wait()
	if found := fetcher.Found(); !found {
		t.Fail()
	} else {
		fmt.Println(fetcher.Result())
	}
}
