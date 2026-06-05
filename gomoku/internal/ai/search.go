package ai

import (
	"gomoku/internal/board"
	"math"
	"sort"
)

func (a *AI) minimaxMove(b *board.Board, p board.Player, depth int) board.Point {
	bestVal := math.MinInt32
	var bestMove board.Point

	// Get candidate moves (only consider points near existing stones to reduce branching)
	candidates := getCandidateMoves(b)
	if len(candidates) == 0 {
		return board.Point{X: board.Size / 2, Y: board.Size / 2} // center
	}

	// Sort candidates by heuristic to improve pruning speed
	sort.Slice(candidates, func(i, j int) bool {
		// Evaluate moves. Combine attack and slightly favor defense in candidates.
		si := EvaluatePosition(b, candidates[i].X, candidates[i].Y, p) + EvaluatePosition(b, candidates[i].X, candidates[i].Y, opponent(p))
		sj := EvaluatePosition(b, candidates[j].X, candidates[j].Y, p) + EvaluatePosition(b, candidates[j].X, candidates[j].Y, opponent(p))
		return si > sj
	})

	// Limit candidates to top 15 to keep search tree manageable
	if len(candidates) > 15 {
		candidates = candidates[:15]
	}

	for _, move := range candidates {
		// Check for immediate win
		if EvaluatePosition(b, move.X, move.Y, p) >= 1000000 {
			return move
		}
		// Check for immediate block
		if EvaluatePosition(b, move.X, move.Y, opponent(p)) >= 1000000 {
			bestMove = move // Remember it, but keep searching just in case we have a win
		}

		// Try move
		b.Grid[move.Y][move.X] = p

		val := minimax(b, depth-1, math.MinInt32, math.MaxInt32, false, p)

		// Undo move
		b.Grid[move.Y][move.X] = board.Empty

		if val > bestVal {
			bestVal = val
			bestMove = move
		}
	}

	if bestMove.X == 0 && bestMove.Y == 0 && b.Grid[0][0] != board.Empty {
		// Fallback if minimax failed to find a valid move
		return a.evaluateMove(b, p)
	}

	return bestMove
}

func minimax(b *board.Board, depth int, alpha, beta int, maximizingPlayer bool, originalPlayer board.Player) int {
	// Terminal node evaluation
	if depth == 0 {
		return evaluateBoard(b, originalPlayer)
	}

	candidates := getCandidateMoves(b)
	if len(candidates) == 0 {
		return 0 // Draw
	}

	// Limit branch factor deeply
	if len(candidates) > 12 {
		candidates = candidates[:12]
	}

	if maximizingPlayer {
		maxEval := math.MinInt32
		for _, move := range candidates {
			b.Grid[move.Y][move.X] = originalPlayer

			// Fast win check during search
			if EvaluatePosition(b, move.X, move.Y, originalPlayer) >= 1000000 {
				b.Grid[move.Y][move.X] = board.Empty
				return 10000000 + depth // Prefer faster wins
			}

			eval := minimax(b, depth-1, alpha, beta, false, originalPlayer)
			b.Grid[move.Y][move.X] = board.Empty

			if eval > maxEval {
				maxEval = eval
			}
			if alpha < eval {
				alpha = eval
			}
			if beta <= alpha {
				break // Beta cut-off
			}
		}
		return maxEval
	} else {
		minEval := math.MaxInt32
		opp := opponent(originalPlayer)
		for _, move := range candidates {
			b.Grid[move.Y][move.X] = opp

			// Fast loss check during search
			if EvaluatePosition(b, move.X, move.Y, opp) >= 1000000 {
				b.Grid[move.Y][move.X] = board.Empty
				return -10000000 - depth // Prefer delayed losses
			}

			eval := minimax(b, depth-1, alpha, beta, true, originalPlayer)
			b.Grid[move.Y][move.X] = board.Empty

			if eval < minEval {
				minEval = eval
			}
			if beta > eval {
				beta = eval
			}
			if beta <= alpha {
				break // Alpha cut-off
			}
		}
		return minEval
	}
}

func getCandidateMoves(b *board.Board) []board.Point {
	var candidates []board.Point
	visited := make(map[board.Point]bool)

	for y := 0; y < board.Size; y++ {
		for x := 0; x < board.Size; x++ {
			if b.Grid[y][x] != board.Empty {
				// Add neighbors within a radius of 2
				for dy := -2; dy <= 2; dy++ {
					for dx := -2; dx <= 2; dx++ {
						nx, ny := x+dx, y+dy
						if nx >= 0 && nx < board.Size && ny >= 0 && ny < board.Size {
							if b.Grid[ny][nx] == board.Empty {
								pt := board.Point{X: nx, Y: ny}
								if !visited[pt] {
									visited[pt] = true
									candidates = append(candidates, pt)
								}
							}
						}
					}
				}
			}
		}
	}
	return candidates
}

func evaluateBoard(b *board.Board, p board.Player) int {
	score := 0
	opp := opponent(p)
	for y := 0; y < board.Size; y++ {
		for x := 0; x < board.Size; x++ {
			if b.Grid[y][x] == p {
				score += EvaluatePosition(b, x, y, p)
			} else if b.Grid[y][x] == opp {
				// Opponent's score counts against us heavily
				score -= EvaluatePosition(b, x, y, opp) * 12 / 10
			}
		}
	}
	return score
}
