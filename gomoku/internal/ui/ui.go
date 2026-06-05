package ui

import (
	"fmt"
	"image/color"
	"math"

	"gomoku/internal/ai"
	"gomoku/internal/board"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type UIState struct {
	boardSize float32
	cellSize  float32
	offsetX   float32
	offsetY   float32
}

func NewUIState() *UIState {
	return &UIState{}
}

func (u *UIState) Update(w, h int) {
	// Calculate responsive board size
	minDim := float32(w)
	if float32(h) < minDim {
		minDim = float32(h)
	}

	u.boardSize = minDim * 0.8
	u.cellSize = u.boardSize / float32(board.Size)
	u.offsetX = (float32(w) - u.boardSize) / 2
	u.offsetY = (float32(h) - u.boardSize) / 2
}

func (u *UIState) ScreenToBoard(sx, sy int, w, h int) (int, int) {
	bx := int(math.Floor(float64((float32(sx) - u.offsetX) / u.cellSize)))
	by := int(math.Floor(float64((float32(sy) - u.offsetY) / u.cellSize)))
	return bx, by
}

func (u *UIState) DrawMenu(screen *ebiten.Image, w, h int, diff ai.Difficulty) {
	cx := float32(w) / 2
	cy := float32(h) / 2

	// Title
	text.Draw(screen, "GOMOKU (五子棋)", titleFont, int(cx)-100, int(cy)-100, color.Black)

	// Start Button
	vector.DrawFilledRect(screen, cx-120, cy-25, 240, 50, color.RGBA{100, 200, 100, 255}, true)
	text.Draw(screen, "START GAME (開始)", normalFont, int(cx)-90, int(cy)+5, color.White)

	// Difficulty Settings
	diffs := []ai.Difficulty{ai.Easy, ai.Normal, ai.Hard}
	names := []string{"EASY(易)", "NORMAL(中)", "HARD(難)"}
	colors := []color.Color{
		color.RGBA{200, 200, 200, 255},
		color.RGBA{200, 200, 200, 255},
		color.RGBA{200, 200, 200, 255},
	}

	for i, d := range diffs {
		if d == diff {
			colors[i] = color.RGBA{100, 150, 250, 255}
		}
		bx := cx - 180 + float32(i*120) // adjusted spacing for longer text
		by := cy + 50
		vector.DrawFilledRect(screen, bx, by, 110, 40, colors[i], true)
		text.Draw(screen, names[i], normalFont, int(bx)+15, int(by)+25, color.Black)
	}
}

func (u *UIState) DrawBoard(screen *ebiten.Image, b *board.Board, w, h int, hoverX, hoverY int, humanTurn bool) {
	// Draw Wood Board Background
	vector.DrawFilledRect(screen, u.offsetX, u.offsetY, u.boardSize, u.boardSize, color.RGBA{222, 184, 135, 255}, true)

	// Draw Grid Lines
	lineColor := color.RGBA{60, 40, 20, 255}
	for i := 0; i < board.Size; i++ {
		// Vertical
		x := u.offsetX + float32(i)*u.cellSize + u.cellSize/2
		vector.StrokeLine(screen, x, u.offsetY+u.cellSize/2, x, u.offsetY+u.boardSize-u.cellSize/2, 2, lineColor, true)

		// Horizontal
		y := u.offsetY + float32(i)*u.cellSize + u.cellSize/2
		vector.StrokeLine(screen, u.offsetX+u.cellSize/2, y, u.offsetX+u.boardSize-u.cellSize/2, y, 2, lineColor, true)
	}

	// Draw Stones
	radius := u.cellSize * 0.4
	for y := 0; y < board.Size; y++ {
		for x := 0; x < board.Size; x++ {
			p := b.Grid[y][x]
			cx := u.offsetX + float32(x)*u.cellSize + u.cellSize/2
			cy := u.offsetY + float32(y)*u.cellSize + u.cellSize/2

			if p != board.Empty {
				c := color.RGBA{0, 0, 0, 255} // Black
				if p == board.White {
					c = color.RGBA{255, 255, 255, 255} // White
				}
				// Draw shadow
				vector.DrawFilledCircle(screen, cx+2, cy+3, radius, color.RGBA{0, 0, 0, 100}, true)
				// Draw stone
				vector.DrawFilledCircle(screen, cx, cy, radius, c, true)

				// Highlight last move
				if b.LastMove.P.X == x && b.LastMove.P.Y == y {
					vector.DrawFilledCircle(screen, cx, cy, radius*0.3, color.RGBA{255, 0, 0, 255}, true)
				}
			}
		}
	}

	// Draw Hover
	if humanTurn && hoverX >= 0 && hoverX < board.Size && hoverY >= 0 && hoverY < board.Size {
		if b.Grid[hoverY][hoverX] == board.Empty {
			cx := u.offsetX + float32(hoverX)*u.cellSize + u.cellSize/2
			cy := u.offsetY + float32(hoverY)*u.cellSize + u.cellSize/2
			vector.DrawFilledCircle(screen, cx, cy, radius, color.RGBA{0, 0, 0, 100}, true)
		}
	}

	// Draw Winning Line
	if len(b.WinningLines) > 0 {
		for _, pt := range b.WinningLines {
			cx := u.offsetX + float32(pt.X)*u.cellSize + u.cellSize/2
			cy := u.offsetY + float32(pt.Y)*u.cellSize + u.cellSize/2
			vector.DrawFilledCircle(screen, cx, cy, radius*0.5, color.RGBA{255, 215, 0, 200}, true)
		}
	}
}

func (u *UIState) DrawHUD(screen *ebiten.Image, turn board.Player, diff ai.Difficulty) {
	text.Draw(screen, fmt.Sprintf("Turn: %s", turn.String()), normalFont, 20, 30, color.Black)
	diffStr := "EASY(易)"
	if diff == ai.Normal {
		diffStr = "NORMAL(中)"
	} else if diff == ai.Hard {
		diffStr = "HARD(難)"
	}
	text.Draw(screen, fmt.Sprintf("AI Diff: %s", diffStr), normalFont, 20, 50, color.Black)
}

func (u *UIState) DrawGameOver(screen *ebiten.Image, winner board.Player, w, h int) {
	cx := float32(w) / 2
	cy := float32(h) / 2

	// Overlay
	vector.DrawFilledRect(screen, 0, 0, float32(w), float32(h), color.RGBA{0, 0, 0, 150}, true)

	msg := "DRAW (平)"
	if winner != board.Empty {
		msg = fmt.Sprintf("%s WINS! (勝!)", winner.String())
	}

	text.Draw(screen, msg, titleFont, int(cx)-100, int(cy)-20, color.White)
	text.Draw(screen, "Click to return (點撃返回)", normalFont, int(cx)-120, int(cy)+30, color.White)
}
