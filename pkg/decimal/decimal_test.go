package decimal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecimalParsing(t *testing.T) {
	d, err := FromString("123.45678901")
	assert.NoError(t, err)
	assert.Equal(t, "123.45678901", d.String())

	_, err = FromString("NaN")
	assert.Error(t, err)

	_, err = FromString("Infinity")
	assert.Error(t, err)

	rf := RoundFiat(MustFromString("123.456"))
	assert.Equal(t, "123.46", rf.String())

	rc := RoundCrypto(MustFromString("1.123456789"))
	assert.Equal(t, "1.12345679", rc.String())
}
