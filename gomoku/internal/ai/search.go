package ai

import (
	"gomoku/internal/board"
	"math"
	"math/rand"
	"sort"
)

func (a *AI) findBestMove(b *board.Board, p board.Player) board.Point {
	opp := opponent(p)
	candidates := getCandidateMoves(b)
	if len(candidates) == 0 {
		return board.Point{X: b.Size / 2, Y: b.Size / 2}
	}

	// 1. Level 1: EASY
	if a.Difficulty == Easy {
		return a.easyMove(b, p, candidates)
	}

	// 2. Strict Tactical Interceptions (Universal for Normal, Hard, Hell)
	// Priority 1: If I can win immediately, PLAY IT!
	for _, m := range candidates {
		if EvaluatePosition(b, m.X, m.Y, p) >= ScoreFive {
			return m
		}
	}

	// Priority 2: If Opponent can win immediately, MUST BLOCK IT!
	for _, m := range candidates {
		if EvaluatePosition(b, m.X, m.Y, opp) >= ScoreFive {
			return m
		}
	}

	// Priority 3: If I can make an Open Four or Double Four, PLAY IT (guaranteed win next turn)
	for _, m := range candidates {
		if EvaluatePosition(b, m.X, m.Y, p) >= ScoreOpenFour {
			return m
		}
	}

	// Priority 4: If Opponent can make an Open Four or Double Four, MUST BLOCK IT!
	for _, m := range candidates {
		if EvaluatePosition(b, m.X, m.Y, opp) >= ScoreOpenFour {
			return m
		}
	}

	// 3. Level 2: NORMAL (Greedy 1-ply tactical scoring)
	if a.Difficulty == Normal {
		return a.greedyMove(b, p, candidates)
	}

	// Priority 5: If I can make Four-Three or Double-Three, PLAY IT!
	for _, m := range candidates {
		if EvaluatePosition(b, m.X, m.Y, p) >= ScoreDoubleThree {
			return m
		}
	}

	// Priority 6: If Opponent can make Four-Three or Double-Three, MUST BLOCK IT!
	for _, m := range candidates {
		if EvaluatePosition(b, m.X, m.Y, opp) >= ScoreDoubleThree {
			return m
		}
	}

	// 4. Level 4: HELL (VCF Forced Kill Search + Deep Search)
	if a.Difficulty == Hell {
		// VCF 1: Can AI force win via continuous checks?
		if win, vcfMove := a.solveVCF(b, p, 10); win {
			return vcfMove
		}
		// VCF 2: Can Opponent force win via continuous checks? Block the trigger move!
		if win, oppVcfMove := a.solveVCF(b, opp, 8); win {
			return oppVcfMove
		}
		// Hell Deep Minimax (depth 4)
		return a.minimaxSearch(b, p, 4, candidates)
	}

	// 5. Level 3: HARD (Minimax depth 3)
	return a.minimaxSearch(b, p, 3, candidates)
}

func (a *AI) easyMove(b *board.Board, p board.Player, candidates []board.Point) board.Point {
	opp := opponent(p)

	// Easy AI has 60% chance to block a 5-in-a-row
	if rand.Float64() < 0.6 {
		for _, m := range candidates {
			if EvaluatePosition(b, m.X, m.Y, p) >= ScoreFive {
				return m
			}
			if EvaluatePosition(b, m.X, m.Y, opp) >= ScoreFive {
				return m
			}
		}
	}

	// 35% chance to pick purely random candidate near stones
	if rand.Float64() < 0.35 && len(candidates) > 0 {
		return candidates[rand.Intn(len(candidates))]
	}

	return a.greedyMove(b, p, candidates)
}

func (a *AI) greedyMove(b *board.Board, p board.Player, candidates []board.Point) board.Point {
	opp := opponent(p)
	bestScore := -1
	var bestMoves []board.Point

	for _, m := range candidates {
		myScore := EvaluatePosition(b, m.X, m.Y, p)
		oppScore := EvaluatePosition(b, m.X, m.Y, opp)
		// Defense weighted 1.25x
		score := myScore + (oppScore * 5 / 4)

		if score > bestScore {
			bestScore = score
			bestMoves = []board.Point{m}
		} else if score == bestScore {
			bestMoves = append(bestMoves, m)
		}
	}

	if len(bestMoves) > 0 {
		return bestMoves[rand.Intn(len(bestMoves))]
	}
	return candidates[0]
}

func (a *AI) minimaxSearch(b *board.Board, p board.Player, depth int, candidates []board.Point) board.Point {
	opp := opponent(p)

	// Sort candidates by heuristic score descending
	sort.Slice(candidates, func(i, j int) bool {
		si := EvaluatePosition(b, candidates[i].X, candidates[i].Y, p) + EvaluatePosition(b, candidates[i].X, candidates[i].Y, opp)*6/5
		sj := EvaluatePosition(b, candidates[j].X, candidates[j].Y, p) + EvaluatePosition(b, candidates[j].X, candidates[j].Y, opp)*6/5
		return si > sj
	})

	limit := 14
	if a.Difficulty == Hell {
		limit = 18
	}
	if len(candidates) > limit {
		candidates = candidates[:limit]
	}

	bestVal := math.MinInt32
	bestMove := candidates[0]
	alpha := math.MinInt32
	beta := math.MaxInt32

	for _, move := range candidates {
		b.Grid[move.Y][move.X] = p
		val := a.minimax(b, depth-1, alpha, beta, false, p)
		b.Grid[move.Y][move.X] = board.Empty

		if val > bestVal {
			bestVal = val
			bestMove = move
		}
		if val > alpha {
			alpha = val
		}
		if beta <= alpha {
			break
		}
	}

	return bestMove
}

func (a *AI) minimax(b *board.Board, depth int, alpha, beta int, maximizingPlayer bool, originalPlayer board.Player) int {
	if depth == 0 {
		return a.evaluateBoardState(b, originalPlayer)
	}

	candidates := getCandidateMoves(b)
	if len(candidates) == 0 {
		return 0
	}

	opp := opponent(originalPlayer)
	limit := 10
	if a.Difficulty == Hell {
		limit = 14
	}

	if maximizingPlayer {
		sort.Slice(candidates, func(i, j int) bool {
			return EvaluatePosition(b, candidates[i].X, candidates[i].Y, originalPlayer) > EvaluatePosition(b, candidates[j].X, candidates[j].Y, originalPlayer)
		})
		if len(candidates) > limit {
			candidates = candidates[:limit]
		}

		maxEval := math.MinInt32
		for _, move := range candidates {
			if EvaluatePosition(b, move.X, move.Y, originalPlayer) >= ScoreFive {
				return ScoreFive + depth*1000
			}

			b.Grid[move.Y][move.X] = originalPlayer
			eval := a.minimax(b, depth-1, alpha, beta, false, originalPlayer)
			b.Grid[move.Y][move.X] = board.Empty

			if eval > maxEval {
				maxEval = eval
			}
			if eval > alpha {
				alpha = eval
			}
			if beta <= alpha {
				break
			}
		}
		return maxEval
	} else {
		sort.Slice(candidates, func(i, j int) bool {
			return EvaluatePosition(b, candidates[i].X, candidates[i].Y, opp) > EvaluatePosition(b, candidates[j].X, candidates[j].Y, opp)
		})
		if len(candidates) > limit {
			candidates = candidates[:limit]
		}

		minEval := math.MaxInt32
		for _, move := range candidates {
			if EvaluatePosition(b, move.X, move.Y, opp) >= ScoreFive {
				return -ScoreFive - depth*1000
			}

			b.Grid[move.Y][move.X] = opp
			eval := a.minimax(b, depth-1, alpha, beta, true, originalPlayer)
			b.Grid[move.Y][move.X] = board.Empty

			if eval < minEval {
				minEval = eval
			}
			if eval < beta {
				beta = eval
			}
			if beta <= alpha {
				break
			}
		}
		return minEval
	}
}

// solveVCF recursively checks if player p has an unstoppable series of continuous four-checks
func (a *AI) solveVCF(b *board.Board, p board.Player, maxDepth int) (bool, board.Point) {
	if maxDepth <= 0 {
		return false, board.Point{-1, -1}
	}

	candidates := getCandidateMoves(b)
	for _, m := range candidates {
		score := EvaluatePosition(b, m.X, m.Y, p)
		if score >= ScoreFive {
			return true, m
		}
		if score >= ScoreOpenFour {
			return true, m
		}

		// Only check moves that produce a four (continuous threat)
		if score >= ScoreBlockedFour {
			b.Grid[m.Y][m.X] = p

			// Opponent MUST block the spot that would create Five for p
			opp := opponent(p)
			var forcedDefense *board.Point
			for _, defPt := range getCandidateMoves(b) {
				if EvaluatePosition(b, defPt.X, defPt.Y, p) >= ScoreFive {
					pt := defPt
					forcedDefense = &pt
					break
				}
			}

			if forcedDefense == nil {
				// No single point can stop p from winning
				b.Grid[m.Y][m.X] = board.Empty
				return true, m
			}

			// Simulate opponent's forced block
			b.Grid[forcedDefense.Y][forcedDefense.X] = opp

			// Does opponent's block accidentally give opponent an immediate win?
			if EvaluatePosition(b, forcedDefense.X, forcedDefense.Y, opp) < ScoreFive {
				// Continue VCF search
				win, _ := a.solveVCF(b, p, maxDepth-1)
				if win {
					b.Grid[forcedDefense.Y][forcedDefense.X] = board.Empty
					b.Grid[m.Y][m.X] = board.Empty
					return true, m
				}
			}

			b.Grid[forcedDefense.Y][forcedDefense.X] = board.Empty
			b.Grid[m.Y][m.X] = board.Empty
		}
	}
	return false, board.Point{-1, -1}
}

func (a *AI) evaluateBoardState(b *board.Board, p board.Player) int {
	myScore := 0
	oppScore := 0
	opp := opponent(p)

	for y := 0; y < b.Size; y++ {
		for x := 0; x < b.Size; x++ {
			if b.Grid[y][x] == p {
				myScore += EvaluatePosition(b, x, y, p)
			} else if b.Grid[y][x] == opp {
				oppScore += EvaluatePosition(b, x, y, opp)
			}
		}
	}
	return myScore - oppScore*6/5
}

func getCandidateMoves(b *board.Board) []board.Point {
	var candidates []board.Point
	visited := make(map[board.Point]bool)
	hasStones := false

	for y := 0; y < b.Size; y++ {
		for x := 0; x < b.Size; x++ {
			if b.Grid[y][x] != board.Empty {
				hasStones = true
				for dy := -2; dy <= 2; dy++ {
					for dx := -2; dx <= 2; dx++ {
						nx, ny := x+dx, y+dy
						if nx >= 0 && nx < b.Size && ny >= 0 && ny < b.Size {
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

	if !hasStones {
		return []board.Point{{X: b.Size / 2, Y: b.Size / 2}}
	}
	return candidates
}
