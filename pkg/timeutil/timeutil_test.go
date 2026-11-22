package timeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeutilConversion(t *testing.T) {
	nowSec := time.Now().Unix()
	s, err := ToSec(nowSec)
	assert.NoError(t, err)
	assert.Equal(t, nowSec, s)

	nowMilli := time.Now().UnixMilli()
	m, err := ToMilli(nowMilli)
	assert.NoError(t, err)
	assert.Equal(t, nowMilli, m)

	tm, err := ToTime(nowSec)
	assert.NoError(t, err)
	assert.Equal(t, nowSec, tm.Unix())
}
