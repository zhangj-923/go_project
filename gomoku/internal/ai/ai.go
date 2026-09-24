package ai

import (
	"gomoku/internal/board"
)

type Difficulty int

const (
	Easy Difficulty = iota
	Normal
	Hard
	Hell
)

type AI struct {
	Difficulty Difficulty
}

func NewAI(diff Difficulty) *AI {
	return &AI{Difficulty: diff}
}

func (a *AI) GetBestMove(b *board.Board, p board.Player) board.Point {
	return a.findBestMove(b, p)
}
