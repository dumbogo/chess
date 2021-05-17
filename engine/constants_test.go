package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToSquareIdentifier(t *testing.T) {
	assert := assert.New(t)
	result, ok := StringToSquareIdentifier("A1")
	assert.Equal(A1, result)
	assert.Equal(true, ok)

	result2, ok := StringToSquareIdentifier("sss")
	assert.Equal(uint8(0), uint8(result2))
	assert.Equal(false, ok)
}
