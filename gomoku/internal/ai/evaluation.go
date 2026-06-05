package ai

import (
	"gomoku/internal/board"
	"math/rand"
)

// Heuristic Evaluation function for Normal/Hard AI
func (a *AI) evaluateMove(b *board.Board, p board.Player) board.Point {
	bestScore := -1
	var bestMoves []board.Point

	for y := 0; y < board.Size; y++ {
		for x := 0; x < board.Size; x++ {
			if b.Grid[y][x] == board.Empty {
				// Normal difficulty only checks 1 step ahead based on scoring
				// Multiply defense score slightly more to prioritize blocking
				score := EvaluatePosition(b, x, y, p) + (EvaluatePosition(b, x, y, opponent(p)) * 12 / 10)

				if score > bestScore {
					bestScore = score
					bestMoves = []board.Point{{X: x, Y: y}}
				} else if score == bestScore {
					bestMoves = append(bestMoves, board.Point{X: x, Y: y})
				}
			}
		}
	}

	if len(bestMoves) > 0 {
		return bestMoves[rand.Intn(len(bestMoves))]
	}
	return a.randomMove(b)
}

func opponent(p board.Player) board.Player {
	if p == board.Black {
		return board.White
	}
	return board.Black
}

// EvaluatePosition returns a score for placing a stone of player p at (x,y)
func EvaluatePosition(b *board.Board, x, y int, p board.Player) int {
	directions := []board.Point{
		{X: 1, Y: 0},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: -1},
	}

	totalScore := 0

	for _, d := range directions {
		consecutive := 1
		blockedEnd1 := false
		blockedEnd2 := false

		// Check Forward
		for i := 1; i <= 4; i++ {
			nx, ny := x+d.X*i, y+d.Y*i
			if nx < 0 || nx >= board.Size || ny < 0 || ny >= board.Size {
				blockedEnd1 = true
				break
			}
			if b.Grid[ny][nx] == p {
				consecutive++
			} else if b.Grid[ny][nx] == board.Empty {
				break
			} else {
				blockedEnd1 = true
				break
			}
		}

		// Check Backward
		for i := 1; i <= 4; i++ {
			nx, ny := x-d.X*i, y-d.Y*i
			if nx < 0 || nx >= board.Size || ny < 0 || ny >= board.Size {
				blockedEnd2 = true
				break
			}
			if b.Grid[ny][nx] == p {
				consecutive++
			} else if b.Grid[ny][nx] == board.Empty {
				break
			} else {
				blockedEnd2 = true
				break
			}
		}

		totalScore += getScore(consecutive, blockedEnd1, blockedEnd2)
	}
	return totalScore
}

func getScore(consecutive int, b1, b2 bool) int {
	if consecutive >= 5 {
		return 1000000 // Win
	}
	if b1 && b2 {
		return 0 // completely blocked
	}

	switch consecutive {
	case 4:
		if b1 || b2 {
			return 10000 // Blocked 4 (冲四)
		}
		return 100000 // Open 4 (活四) - Guaranteed win next turn
	case 3:
		if b1 || b2 {
			return 1000 // Blocked 3 (眠三)
		}
		return 10000 // Open 3 (活三)
	case 2:
		if b1 || b2 {
			return 100 // Blocked 2 (眠二)
		}
		return 1000 // Open 2 (活二)
	case 1:
		if b1 || b2 {
			return 10
		}
		return 100 // Open 1
	}
	return 0
}
