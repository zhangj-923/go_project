package ui

import (
	"fmt"
	"image/color"
	"math"

	"gomoku/internal/ai"
	"gomoku/internal/board"
	"gomoku/internal/i18n"

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

func (u *UIState) Update(w, h int, boardGridSize int, isOnlineDesktop bool) {
	if boardGridSize <= 0 {
		boardGridSize = board.DefaultSize
	}
	fw := float32(w)
	fh := float32(h)

	if fw < fh {
		// 手机竖屏自适应：棋盘占宽 92%，留出上下控制区
		u.boardSize = fw * 0.92
		if u.boardSize > fh*0.62 {
			u.boardSize = fh * 0.62
		}
		u.cellSize = u.boardSize / float32(boardGridSize)
		u.offsetX = (fw - u.boardSize) / 2
		u.offsetY = (fh-u.boardSize)/2 - 20
	} else {
		// 宽屏/桌面端：若处于在线宽屏对决，扣除右侧 350px 控制台空间，防止 1280x800 等笔记本重叠
		availW := fw
		if isOnlineDesktop && fw > 960 {
			availW = fw - 350
		}
		minDim := availW
		if fh < minDim {
			minDim = fh
		}
		u.boardSize = minDim * 0.82
		u.cellSize = u.boardSize / float32(boardGridSize)
		u.offsetX = (availW - u.boardSize) / 2
		u.offsetY = (fh - u.boardSize) / 2
	}
}

func (u *UIState) ScreenToBoard(sx, sy int, w, h int, boardGridSize int) (int, int) {
	bx := int(math.Floor(float64((float32(sx) - u.offsetX) / u.cellSize)))
	by := int(math.Floor(float64((float32(sy) - u.offsetY) / u.cellSize)))
	return bx, by
}

func (u *UIState) DrawMenu(screen *ebiten.Image, w, h int, isAI bool, diff ai.Difficulty, bSize int, lang i18n.Language) {
	cx := float32(w) / 2
	cy := float32(h) / 2

	// 1. Language Toggle Button (Top-Right)
	langBtnW := float32(130)
	langBtnH := float32(36)
	langBtnX := float32(w) - langBtnW - 20
	langBtnY := float32(20)
	vector.DrawFilledRect(screen, langBtnX, langBtnY, langBtnW, langBtnH, color.RGBA{45, 52, 54, 230}, true)
	vector.StrokeRect(screen, langBtnX, langBtnY, langBtnW, langBtnH, 1, color.RGBA{178, 190, 195, 200}, true)
	text.Draw(screen, i18n.T(lang, i18n.LangName), normalFont, int(langBtnX)+15, int(langBtnY)+24, color.White)

	// 2. Title & Subtitle
	titleStr := i18n.T(lang, i18n.Title)
	text.Draw(screen, titleStr, titleFont, int(cx)-130, int(cy)-175, color.RGBA{45, 35, 25, 255})
	subStr := i18n.T(lang, i18n.Subtitle)
	text.Draw(screen, subStr, normalFont, int(cx)-75, int(cy)-140, color.RGBA{110, 95, 80, 255})

	// 3. Start Game Button
	startW := float32(260)
	startH := float32(52)
	startX := cx - startW/2
	startY := cy - 110
	vector.DrawFilledRect(screen, startX, startY, startW, startH, color.RGBA{46, 170, 90, 255}, true)
	vector.StrokeRect(screen, startX, startY, startW, startH, 2, color.RGBA{35, 135, 70, 255}, true)
	startText := i18n.T(lang, i18n.StartGame)
	text.Draw(screen, startText, normalFont, int(startX)+70, int(startY)+33, color.White)

	// 4. Mode Selection Row
	modeY := cy - 35
	text.Draw(screen, i18n.T(lang, i18n.ModeSelect)+":", normalFont, int(cx)-195, int(modeY)+24, color.RGBA{60, 50, 40, 255})
	modeBtns := []string{i18n.T(lang, i18n.ModeVsAI), i18n.T(lang, i18n.ModePvP)}
	for i, name := range modeBtns {
		btnX := cx - 115 + float32(i*155)
		c := color.RGBA{210, 200, 185, 255}
		tc := color.Black
		if (i == 0 && isAI) || (i == 1 && !isAI) {
			c = color.RGBA{70, 130, 240, 255}
			tc = color.White
		}
		vector.DrawFilledRect(screen, btnX, modeY, 140, 36, c, true)
		text.Draw(screen, name, normalFont, int(btnX)+20, int(modeY)+24, tc)
	}

	// 5. Board Size Selection Row
	sizeY := cy + 20
	text.Draw(screen, i18n.T(lang, i18n.BoardSize)+":", normalFont, int(cx)-195, int(sizeY)+24, color.RGBA{60, 50, 40, 255})
	sizes := []int{11, 15, 19}
	sizeNames := []string{i18n.T(lang, i18n.Size11), i18n.T(lang, i18n.Size15), i18n.T(lang, i18n.Size19)}
	for i, s := range sizes {
		btnX := cx - 95 + float32(i*105)
		c := color.RGBA{210, 200, 185, 255}
		tc := color.Black
		if s == bSize {
			c = color.RGBA{225, 140, 40, 255}
			tc = color.White
		}
		vector.DrawFilledRect(screen, btnX, sizeY, 95, 36, c, true)
		text.Draw(screen, sizeNames[i], smallFont, int(btnX)+8, int(sizeY)+23, tc)
	}

	// 6. Difficulty Selection Row (Only for Vs AI)
	if isAI {
		diffY := cy + 75
		text.Draw(screen, i18n.T(lang, i18n.Difficulty)+":", normalFont, int(cx)-195, int(diffY)+24, color.RGBA{60, 50, 40, 255})

		diffs := []ai.Difficulty{ai.Easy, ai.Normal, ai.Hard, ai.Hell}
		diffNames := []string{
			i18n.T(lang, i18n.DiffEasy),
			i18n.T(lang, i18n.DiffNormal),
			i18n.T(lang, i18n.DiffHard),
			i18n.T(lang, i18n.DiffHell),
		}

		for i, d := range diffs {
			btnX := cx - 95 + float32(i*78)
			c := color.RGBA{210, 200, 185, 255}
			tc := color.Black
			if d == diff {
				if d == ai.Hell {
					c = color.RGBA{220, 45, 45, 255} // Hell mode prominent crimson red
				} else {
					c = color.RGBA{70, 130, 240, 255}
				}
				tc = color.White
			}
			vector.DrawFilledRect(screen, btnX, diffY, 70, 36, c, true)
			text.Draw(screen, diffNames[i], smallFont, int(btnX)+16, int(diffY)+23, tc)
		}
	}
}

func (u *UIState) DrawBoard(screen *ebiten.Image, b *board.Board, w, h int, hoverX, hoverY int, humanTurn bool) {
	// Draw Wood Board Background
	vector.DrawFilledRect(screen, u.offsetX, u.offsetY, u.boardSize, u.boardSize, color.RGBA{225, 185, 130, 255}, true)
	// Outer border frame
	vector.StrokeRect(screen, u.offsetX, u.offsetY, u.boardSize, u.boardSize, 4, color.RGBA{115, 80, 40, 255}, true)

	// Draw Grid Lines
	lineColor := color.RGBA{70, 45, 20, 230}
	for i := 0; i < b.Size; i++ {
		// Vertical
		x := u.offsetX + float32(i)*u.cellSize + u.cellSize/2
		vector.StrokeLine(screen, x, u.offsetY+u.cellSize/2, x, u.offsetY+u.boardSize-u.cellSize/2, 1.5, lineColor, true)

		// Horizontal
		y := u.offsetY + float32(i)*u.cellSize + u.cellSize/2
		vector.StrokeLine(screen, u.offsetX+u.cellSize/2, y, u.offsetX+u.boardSize-u.cellSize/2, y, 1.5, lineColor, true)
	}

	// Draw Center Point and Star Points
	starPoints := getStarPoints(b.Size)
	for _, pt := range starPoints {
		cx := u.offsetX + float32(pt.X)*u.cellSize + u.cellSize/2
		cy := u.offsetY + float32(pt.Y)*u.cellSize + u.cellSize/2
		vector.DrawFilledCircle(screen, cx, cy, 3.5, lineColor, true)
	}

	// Draw Stones
	radius := u.cellSize * 0.42
	for y := 0; y < b.Size; y++ {
		for x := 0; x < b.Size; x++ {
			p := b.Grid[y][x]
			cx := u.offsetX + float32(x)*u.cellSize + u.cellSize/2
			cy := u.offsetY + float32(y)*u.cellSize + u.cellSize/2

			if p != board.Empty {
				c := color.RGBA{25, 25, 25, 255} // Black
				if p == board.White {
					c = color.RGBA{245, 245, 245, 255} // White
				}
				// Draw shadow
				vector.DrawFilledCircle(screen, cx+2, cy+3, radius, color.RGBA{0, 0, 0, 70}, true)
				// Draw stone body
				vector.DrawFilledCircle(screen, cx, cy, radius, c, true)
				// Highlight circle for 3D depth
				if p == board.White {
					vector.StrokeCircle(screen, cx, cy, radius, 1, color.RGBA{200, 200, 200, 255}, true)
				}

				// Highlight last move with a red indicator
				if b.LastMove.P.X == x && b.LastMove.P.Y == y {
					vector.DrawFilledCircle(screen, cx, cy, radius*0.35, color.RGBA{235, 50, 50, 255}, true)
				}
			}
		}
	}

	// Draw Hover Ghost Stone
	if humanTurn && hoverX >= 0 && hoverX < b.Size && hoverY >= 0 && hoverY < b.Size {
		if b.Grid[hoverY][hoverX] == board.Empty {
			cx := u.offsetX + float32(hoverX)*u.cellSize + u.cellSize/2
			cy := u.offsetY + float32(hoverY)*u.cellSize + u.cellSize/2
			vector.DrawFilledCircle(screen, cx, cy, radius, color.RGBA{50, 120, 230, 110}, true)
		}
	}

	// Draw Winning Stones Highlight
	if len(b.WinningLines) > 0 {
		for _, pt := range b.WinningLines {
			cx := u.offsetX + float32(pt.X)*u.cellSize + u.cellSize/2
			cy := u.offsetY + float32(pt.Y)*u.cellSize + u.cellSize/2
			vector.DrawFilledCircle(screen, cx, cy, radius*0.5, color.RGBA{255, 215, 0, 210}, true)
		}
	}
}

func getStarPoints(size int) []board.Point {
	if size == 15 {
		return []board.Point{
			{3, 3}, {11, 3}, {7, 7}, {3, 11}, {11, 11},
		}
	} else if size == 19 {
		return []board.Point{
			{3, 3}, {9, 3}, {15, 3},
			{3, 9}, {9, 9}, {15, 9},
			{3, 15}, {9, 15}, {15, 15},
		}
	} else if size == 11 {
		return []board.Point{
			{5, 5},
		}
	}
	return nil
}

func (u *UIState) DrawHUD(screen *ebiten.Image, isAI bool, turn board.Player, diff ai.Difficulty, bSize int, lang i18n.Language) {
	// Top status bar
	var turnText string
	if isAI {
		if turn == board.Black {
			turnText = i18n.T(lang, i18n.TurnBlackPlayer)
		} else {
			turnText = i18n.T(lang, i18n.TurnWhiteAI)
		}
	} else {
		if turn == board.Black {
			turnText = i18n.T(lang, i18n.TurnBlack)
		} else {
			turnText = i18n.T(lang, i18n.TurnWhite)
		}
	}
	text.Draw(screen, turnText, normalFont, 24, 30, color.RGBA{40, 30, 20, 255})

	// Game config badge
	configText := fmt.Sprintf("%dx%d", bSize, bSize)
	if isAI {
		diffStr := i18n.T(lang, i18n.DiffEasy)
		switch diff {
		case ai.Normal:
			diffStr = i18n.T(lang, i18n.DiffNormal)
		case ai.Hard:
			diffStr = i18n.T(lang, i18n.DiffHard)
		case ai.Hell:
			diffStr = i18n.T(lang, i18n.DiffHell)
		}
		configText += " | " + diffStr
	} else {
		configText += " | " + i18n.T(lang, i18n.ModePvP)
	}
	text.Draw(screen, configText, smallFont, 24, 52, color.RGBA{110, 95, 80, 255})
}

func (u *UIState) DrawGameOver(screen *ebiten.Image, winner board.Player, w, h int, isOnline bool, lang i18n.Language) {
	cx := float32(w) / 2
	cy := float32(h) / 2

	// Semi-transparent overlay
	vector.DrawFilledRect(screen, 0, 0, float32(w), float32(h), color.RGBA{0, 0, 0, 160}, true)

	var msg string
	if winner == board.Black {
		msg = i18n.T(lang, i18n.WinBlack)
	} else if winner == board.White {
		msg = i18n.T(lang, i18n.WinWhite)
	} else {
		msg = i18n.T(lang, i18n.DrawGame)
	}

	text.Draw(screen, msg, titleFont, int(cx)-110, int(cy)-20, color.White)
	var returnText string
	offsetX := 95
	if isOnline {
		returnText = i18n.T(lang, i18n.OnlineGameOverHint)
		offsetX = 145
	} else {
		returnText = i18n.T(lang, i18n.ReturnToMenu)
	}
	text.Draw(screen, returnText, normalFont, int(cx)-offsetX, int(cy)+35, color.RGBA{220, 220, 220, 255})
}
