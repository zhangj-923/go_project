package ai

import (
	"gomoku/internal/board"
	"math"
)

// Scores for standard Gomoku tactical shapes
const (
	ScoreFive         = 10000000 // 连五 (必胜)
	ScoreOpenFour     = 1000000  // 活四 (必杀)
	ScoreDoubleFour   = 1000000  // 双冲四 (必杀)
	ScoreFourThree    = 500000   // 冲四活三 (绝杀)
	ScoreDoubleThree  = 200000   // 双活三 (绝杀)
	ScoreBlockedFour  = 100000   // 冲四
	ScoreOpenThree    = 50000    // 活三
	ScoreBlockedThree = 5000     // 眠三
	ScoreOpenTwo      = 2000     // 活二
	ScoreBlockedTwo   = 200      // 眠二
)

// EvaluatePosition evaluates the tactical score of placing player p's stone at (x, y)
func EvaluatePosition(b *board.Board, x, y int, p board.Player) int {
	if b.Grid[y][x] != board.Empty {
		return 0
	}

	directions := [][2]int{
		{1, 0},  // 水平 -
		{0, 1},  // 垂直 |
		{1, 1},  // 正斜 \
		{1, -1}, // 反斜 /
	}

	fives := 0
	openFours := 0
	blockedFours := 0
	openThrees := 0
	blockedThrees := 0
	openTwos := 0

	opp := opponent(p)

	for _, d := range directions {
		dx, dy := d[0], d[1]

		// Build a 9-cell window centered at (x, y)
		// 0: empty, 1: player p, 2: blocked (opponent or wall)
		var line [9]int
		for i := -4; i <= 4; i++ {
			if i == 0 {
				line[4] = 1 // the stone we are placing
				continue
			}
			nx, ny := x+dx*i, y+dy*i
			if nx < 0 || nx >= b.Size || ny < 0 || ny >= b.Size {
				line[i+4] = 2 // Wall
			} else if b.Grid[ny][nx] == opp {
				line[i+4] = 2 // Opponent
			} else if b.Grid[ny][nx] == p {
				line[i+4] = 1 // Player
			} else {
				line[i+4] = 0 // Empty
			}
		}

		// Analyze 5-cell windows in the 9-cell array
		dirMax := 0
		hasOpenFour := false
		hasBlockedFour := false
		hasOpenThree := false
		hasBlockedThree := false
		hasOpenTwo := false

		for s := 0; s <= 4; s++ {
			pCount := 0
			hasBlock := false
			for k := 0; k < 5; k++ {
				if line[s+k] == 2 {
					hasBlock = true
					break
				}
				if line[s+k] == 1 {
					pCount++
				}
			}
			if hasBlock {
				continue
			}

			if pCount == 5 {
				if dirMax < ScoreFive {
					dirMax = ScoreFive
				}
			} else if pCount == 4 {
				leftOpen := (s > 0 && line[s-1] == 0)
				rightOpen := (s+5 < 9 && line[s+5] == 0)
				if leftOpen && rightOpen {
					hasOpenFour = true
				} else if leftOpen || rightOpen {
					hasBlockedFour = true
				}
			} else if pCount == 3 {
				leftOpen := (s > 0 && line[s-1] == 0)
				rightOpen := (s+5 < 9 && line[s+5] == 0)
				if leftOpen && rightOpen {
					hasOpenThree = true
				} else if leftOpen || rightOpen {
					hasBlockedThree = true
				}
			} else if pCount == 2 {
				leftOpen := (s > 0 && line[s-1] == 0)
				rightOpen := (s+5 < 9 && line[s+5] == 0)
				if leftOpen && rightOpen {
					hasOpenTwo = true
				}
			}
		}

		if dirMax == ScoreFive {
			fives++
		} else if hasOpenFour {
			openFours++
		} else if hasBlockedFour {
			blockedFours++
		} else if hasOpenThree {
			openThrees++
		} else if hasBlockedThree {
			blockedThrees++
		} else if hasOpenTwo {
			openTwos++
		}
	}

	// Combine combination threats
	if fives > 0 {
		return ScoreFive
	}
	if openFours > 0 || blockedFours >= 2 {
		return ScoreOpenFour
	}
	if blockedFours > 0 && openThrees > 0 {
		return ScoreFourThree
	}
	if openThrees >= 2 {
		return ScoreDoubleThree
	}

	total := 0
	if blockedFours > 0 {
		total += ScoreBlockedFour * blockedFours
	}
	if openThrees > 0 {
		total += ScoreOpenThree * openThrees
	}
	if blockedThrees > 0 {
		total += ScoreBlockedThree * blockedThrees
	}
	if openTwos > 0 {
		total += ScoreOpenTwo * openTwos
	}

	// Positional center bias (closer to center of board is slightly better)
	center := float64(b.Size-1) / 2.0
	dist := math.Abs(float64(x)-center) + math.Abs(float64(y)-center)
	centerBonus := int((float64(b.Size) - dist) * 5)
	if centerBonus > 0 {
		total += centerBonus
	}

	return total
}

func opponent(p board.Player) board.Player {
	if p == board.Black {
		return board.White
	}
	return board.Black
}
