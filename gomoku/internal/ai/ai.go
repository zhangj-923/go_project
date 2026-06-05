package ai

import (
	"gomoku/internal/board"
	"math/rand"
)

type Difficulty int

const (
	Easy Difficulty = iota
	Normal
	Hard
)

type AI struct {
	Difficulty Difficulty
}

func NewAI(diff Difficulty) *AI {
	return &AI{Difficulty: diff}
}

func (a *AI) GetBestMove(b *board.Board, p board.Player) board.Point {
	switch a.Difficulty {
	case Easy:
		return a.randomMove(b)
	case Normal:
		return a.evaluateMove(b, p) // Simple heuristic
	case Hard:
		return a.minimaxMove(b, p, 3) // Minimax with depth 3
	default:
		return a.randomMove(b)
	}
}

func (a *AI) randomMove(b *board.Board) board.Point {
	var emptyPoints []board.Point
	for y := 0; y < board.Size; y++ {
		for x := 0; x < board.Size; x++ {
			if b.Grid[y][x] == board.Empty {
				emptyPoints = append(emptyPoints, board.Point{X: x, Y: y})
			}
		}
	}
	if len(emptyPoints) == 0 {
		return board.Point{X: -1, Y: -1}
	}
	return emptyPoints[rand.Intn(len(emptyPoints))]
}
