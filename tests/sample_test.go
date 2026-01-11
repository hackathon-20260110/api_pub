package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSample(t *testing.T) {
	// 1=1が成り立つことを確認
	assert.Equal(t, 1, 1)

	// 1=2が成り立たないことを確認
	assert.NotEqual(t, 1, 2)
}
