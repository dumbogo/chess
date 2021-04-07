package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPawn(t *testing.T) {
	assert := assert.New(t)
	pawn := NewPawn(BlackColor)
	assert.Equal(PawnIdentifier, pawn.Identifier())
	assert.Equal(BlackColor, pawn.Color())
}

// TODO: add tests CanMove interface method for each piece
