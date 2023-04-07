package testhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertEqualErrors(t *testing.T, want, got error) {
	if want == nil {
		assert.Nil(t, got)
		return
	}
	if got == nil {
		assert.Nil(t, want)
		return
	}

	assert.Error(t, got)
	assert.Equal(t, want.Error(), got.Error())
}
