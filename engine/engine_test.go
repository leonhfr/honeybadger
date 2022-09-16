package engine

import (
	"testing"

	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	name, author := "NAME", "AUTHOR"
	e := New(name, author)
	assert.Equal(t, name, e.name)
	assert.Equal(t, author, e.author)
	assert.Equal(t, chess.StartingPosition().String(), e.game.Position().String())
}

func TestSetPositionValid(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	e := New("", "")
	err := e.SetPosition(fen)
	if assert.NoError(t, err) {
		assert.Equal(t, fen, e.game.Position().String())
	}
}

func TestSetPositionInvalid(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/BROKEN_ b KQkq e3 0 1"
	e := New("", "")
	err := e.SetPosition(fen)
	assert.Error(t, err)
}

func TestMoveValid(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	e := New("", "")
	err := e.Move("e2e4")
	if assert.NoError(t, err) {
		assert.Equal(t, fen, e.game.Position().String())
	}
}

func TestMoveInvalid(t *testing.T) {
	e := New("", "")
	err := e.Move("e2e5")
	assert.Error(t, err)
}

func TestResetPosition(t *testing.T) {
	e := New("", "")
	e.ResetPosition()
	assert.Equal(t, chess.StartingPosition().String(), e.game.Position().String())
}
