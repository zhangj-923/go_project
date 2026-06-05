package game

// CalculateMove paths the piece given a dice roll and returns the final position and events.
func CalculateMove(board *BoardState, piece *Piece, diceValue int) (int, []MoveEvent) {
	if piece.Position == -1 {
		// Launch
		if diceValue == 6 {
			startPos := StartPositions[piece.Color]
			events := []MoveEvent{{Type: EventMove, Position: startPos}}
			// After launch, we don't jump or fly immediately in some rules, but let's check for stack/eat
			finalPos, eatEvents := checkEat(board, piece, startPos)
			events = append(events, eatEvents...)
			return finalPos, events
		}
		return -1, nil // Cannot move
	}

	var events []MoveEvent
	pos := piece.Position

	// Step by step calculation to handle turn home and bounce
	for i := 0; i < diceValue; i++ {
		// Check if turning home
		if pos == TurnHomePositions[piece.Color] {
			pos = HomeStretchBases[piece.Color]
		} else if pos >= 100 {
			// In home stretch
			homeWin := HomeStretchBases[piece.Color] + 5
			if pos == homeWin {
				// Bounce back!
				pos--
			} else {
				// Check if moving towards home or bouncing
				// To properly bounce, we need to know direction.
				// The easiest way to handle bounce is to calculate remaining steps if we reach win.
				remaining := diceValue - i
				if pos+remaining > homeWin {
					// Overshot
					overshot := (pos + remaining) - homeWin
					pos = homeWin - overshot
					break
				} else {
					pos += remaining
					break
				}
			}
		} else {
			pos = (pos + 1) % 52
		}
	}

	events = append(events, MoveEvent{Type: EventMove, Position: pos})

	if pos >= 100 {
		if pos == HomeStretchBases[piece.Color]+5 {
			events = append(events, MoveEvent{Type: EventHome, Position: pos})
		}
		return pos, events
	}

	// Check if lands on Fly Cell
	isFly, flyTarget := IsFlyCell(piece.Color, pos)
	if isFly {
		events = append(events, MoveEvent{Type: EventFly, Position: flyTarget})
		pos = flyTarget

		// After fly, there's often a jump if it lands on same color.
		// Fly target is always same color.
		// For example, fly from 16 to 28 (Red). 28 is Red. It should jump to 32.
		// Is 28 a turn home? For Red, turn home is 50, so no.
		// Wait, jump rule: land on same color, jump to next same color (+4).
		nextPos := (pos + 4) % 52
		events = append(events, MoveEvent{Type: EventJump, Position: nextPos})
		pos = nextPos
	} else if CellColor(pos) == piece.Color {
		// Normal jump
		nextPos := (pos + 4) % 52
		// If jump lands on fly cell?
		isFlyAfterJump, flyTargetAfterJump := IsFlyCell(piece.Color, nextPos)
		events = append(events, MoveEvent{Type: EventJump, Position: nextPos})
		pos = nextPos

		if isFlyAfterJump {
			events = append(events, MoveEvent{Type: EventFly, Position: flyTargetAfterJump})
			pos = flyTargetAfterJump
		}
	}

	// Check Eat
	finalPos, eatEvents := checkEat(board, piece, pos)
	events = append(events, eatEvents...)

	return finalPos, events
}

func checkEat(board *BoardState, movingPiece *Piece, pos int) (int, []MoveEvent) {
	var events []MoveEvent
	if pos >= 100 || pos == -1 {
		return pos, events // Cannot eat in home stretch or base
	}

	// Find enemies at this pos
	for i := range board.Players {
		p := &board.Players[i]
		if p.Color == movingPiece.Color {
			continue // Ally or self
		}
		for j := range p.Pieces {
			enemyPiece := &p.Pieces[j]
			if enemyPiece.Position == pos {
				// EAT!
				events = append(events, MoveEvent{
					Type:     EventEat,
					Position: -1,
					Target:   enemyPiece.ID, // Need to specify which player's piece
					// Wait, Target int is just PieceID. We need color too. Let's use a struct or just infer.
					// Since there can be stacked pieces, we might eat multiple!
				})
			}
		}
	}
	return pos, events
}

func GetMovablePieces(board *BoardState, color Color, diceValue int) []int {
	var movable []int
	var player *Player
	for i := range board.Players {
		if board.Players[i].Color == color {
			player = &board.Players[i]
			break
		}
	}
	if player == nil {
		return movable
	}

	for _, piece := range player.Pieces {
		if piece.Position == -1 {
			if diceValue == 6 {
				movable = append(movable, piece.ID)
			}
		} else if !IsWinPos(color, piece.Position) {
			movable = append(movable, piece.ID)
		}
	}
	return movable
}
