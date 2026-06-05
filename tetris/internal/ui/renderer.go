package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"tetris/internal/config"
	"tetris/internal/game"
)

var (
	styleBorder   = tcell.StyleDefault.Foreground(tcell.NewRGBColor(100, 100, 120))
	styleText     = tcell.StyleDefault.Foreground(tcell.NewRGBColor(200, 200, 220)).Bold(true)
	styleTitle    = tcell.StyleDefault.Foreground(tcell.NewRGBColor(0, 200, 255)).Bold(true)
	styleScore    = tcell.StyleDefault.Foreground(tcell.NewRGBColor(255, 200, 50)).Bold(true)
	styleGameOver = tcell.StyleDefault.Foreground(tcell.NewRGBColor(255, 60, 60)).Bold(true)
	stylePause    = tcell.StyleDefault.Foreground(tcell.NewRGBColor(255, 200, 0)).Bold(true)
	styleBg       = tcell.StyleDefault.Background(tcell.NewRGBColor(15, 15, 25))
)

type Renderer struct {
	Screen       tcell.Screen
	Config       *config.Config
	BoardOffsetX int
	BoardOffsetY int
}

func NewRenderer(screen tcell.Screen, cfg *config.Config) *Renderer {
	r := &Renderer{
		Screen: screen,
		Config: cfg,
	}
	r.UpdateLayout()
	return r
}

func (r *Renderer) UpdateLayout() {
	sw, sh := r.Screen.Size()
	boardPixelWidth := r.Config.BoardWidth * r.Config.CellWidth
	boardPixelHeight := r.Config.BoardHeight * r.Config.CellHeight

	r.BoardOffsetX = (sw - boardPixelWidth) / 2
	r.BoardOffsetY = (sh - boardPixelHeight) / 2
	if r.BoardOffsetX < 0 {
		r.BoardOffsetX = 0
	}
	if r.BoardOffsetY < 0 {
		r.BoardOffsetY = 0
	}
}

func (r *Renderer) Draw(engine *game.Engine) {
	r.Screen.Clear()
	sw, sh := r.Screen.Size()
	// 填充背景
	for y := 0; y < sh; y++ {
		for x := 0; x < sw; x++ {
			r.Screen.SetContent(x, y, ' ', nil, styleBg)
		}
	}
	r.drawBoard(engine)
	r.drawCurrentPiece(engine)
	r.drawGhost(engine)
	r.drawSidebar(engine)
	r.drawControls()
	if engine.GameOver {
		r.drawGameOver()
	}
	if engine.Paused {
		r.drawPaused()
	}
	r.Screen.Show()
}

func (r *Renderer) drawBoard(engine *game.Engine) {
	ox, oy := r.BoardOffsetX, r.BoardOffsetY
	cw, ch := r.Config.CellWidth, r.Config.CellHeight

	// 绘制棋盘内部
	for row := 0; row < r.Config.BoardHeight; row++ {
		for col := 0; col < r.Config.BoardWidth; col++ {
			x := ox + col*cw
			y := oy + row*ch
			if engine.Board[row][col].Filled {
				st := tcell.StyleDefault.Background(engine.Board[row][col].Color)
				r.fillBlock(x, y, cw, ch, ' ', st)
			} else {
				// 背景交替色
				var st tcell.Style
				if (row+col)%2 == 0 {
					st = tcell.StyleDefault.Background(tcell.NewRGBColor(22, 22, 38))
				} else {
					st = tcell.StyleDefault.Background(tcell.NewRGBColor(18, 18, 32))
				}
				r.fillBlock(x, y, cw, ch, ' ', st)
			}
		}
	}

	boardW := r.Config.BoardWidth * cw
	boardH := r.Config.BoardHeight * ch

	// 绘制边框
	for c := -1; c <= boardW; c++ {
		r.Screen.SetContent(ox+c, oy-1, '═', nil, styleBorder)
		r.Screen.SetContent(ox+c, oy+boardH, '═', nil, styleBorder)
	}
	for row := 0; row < boardH; row++ {
		r.Screen.SetContent(ox-1, oy+row, '║', nil, styleBorder)
		r.Screen.SetContent(ox+boardW, oy+row, '║', nil, styleBorder)
	}

	r.Screen.SetContent(ox-1, oy-1, '╔', nil, styleBorder)
	r.Screen.SetContent(ox+boardW, oy-1, '╗', nil, styleBorder)
	r.Screen.SetContent(ox-1, oy+boardH, '╚', nil, styleBorder)
	r.Screen.SetContent(ox+boardW, oy+boardH, '╝', nil, styleBorder)
}

func (r *Renderer) fillBlock(x, y, w, h int, ch rune, style tcell.Style) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			r.Screen.SetContent(x+dx, y+dy, ch, nil, style)
		}
	}
}

func (r *Renderer) drawCurrentPiece(engine *game.Engine) {
	shape := game.Tetrominoes[engine.Current.TetrominoIdx].Rotations[engine.Current.Rotation]
	col := game.Tetrominoes[engine.Current.TetrominoIdx].Color
	ox, oy := r.BoardOffsetX, r.BoardOffsetY
	cw, ch := r.Config.CellWidth, r.Config.CellHeight

	for row := 0; row < 4; row++ {
		for c := 0; c < 4; c++ {
			if shape[row][c] == 0 {
				continue
			}
			py := engine.Current.Y + row
			px := engine.Current.X + c
			if py >= 0 && py < r.Config.BoardHeight && px >= 0 && px < r.Config.BoardWidth {
				r.fillBlock(ox+px*cw, oy+py*ch, cw, ch, ' ', tcell.StyleDefault.Background(col))
			}
		}
	}
}

func (r *Renderer) drawGhost(engine *game.Engine) {
	if engine.GhostY == engine.Current.Y {
		return
	}
	shape := game.Tetrominoes[engine.Current.TetrominoIdx].Rotations[engine.Current.Rotation]
	ox, oy := r.BoardOffsetX, r.BoardOffsetY
	cw, ch := r.Config.CellWidth, r.Config.CellHeight
	ghostColor := tcell.NewRGBColor(60, 60, 80)
	st := tcell.StyleDefault.Background(ghostColor)

	for row := 0; row < 4; row++ {
		for c := 0; c < 4; c++ {
			if shape[row][c] == 0 {
				continue
			}
			py := engine.GhostY + row
			px := engine.Current.X + c
			if py >= 0 && py < r.Config.BoardHeight && px >= 0 && px < r.Config.BoardWidth {
				r.fillBlock(ox+px*cw, oy+py*ch, cw, ch, '░', st)
			}
		}
	}
}

// 修改后的 drawString 能够正确处理包含全角中文字符的字符串
func drawString(s tcell.Screen, x, y int, style tcell.Style, text string) {
	col := x
	for _, ch := range text {
		s.SetContent(col, y, ch, nil, style)
		col += runewidth.RuneWidth(ch)
	}
}

func (r *Renderer) drawSidebar(engine *game.Engine) {
	ox := r.BoardOffsetX + r.Config.BoardWidth*r.Config.CellWidth + 3
	oy := r.BoardOffsetY

	drawString(r.Screen, ox, oy, styleTitle, "╔══════════════╗")
	drawString(r.Screen, ox, oy+1, styleTitle, "║  俄罗斯方块  ║")
	drawString(r.Screen, ox, oy+2, styleTitle, "╚══════════════╝")

	drawString(r.Screen, ox, oy+4, styleText, "┌── 分数 ───┐")
	drawString(r.Screen, ox, oy+5, styleScore, fmt.Sprintf("   %08d   ", engine.Score))
	drawString(r.Screen, ox, oy+6, styleText, "└───────────┘")

	drawString(r.Screen, ox, oy+8, styleText, fmt.Sprintf("  等级: %d", engine.Level))
	drawString(r.Screen, ox, oy+9, styleText, fmt.Sprintf("  行数: %d", engine.Lines))

	drawString(r.Screen, ox, oy+11, styleText, "┌── 下一个 ─┐")
	nextShape := game.Tetrominoes[engine.Next].Rotations[0]
	nextColor := game.Tetrominoes[engine.Next].Color
	cw, ch := r.Config.CellWidth, r.Config.CellHeight
	for row := 0; row < 4; row++ {
		for c := 0; c < 4; c++ {
			if nextShape[row][c] == 1 {
				r.fillBlock(ox+2+c*cw, oy+12+row*ch, cw, ch, ' ', tcell.StyleDefault.Background(nextColor))
			}
		}
	}
	// 动态调整 Next 框底部高度
	drawString(r.Screen, ox, oy+13+4*ch, styleText, "└───────────┘")
}

func (r *Renderer) drawControls() {
	ox := r.BoardOffsetX - 20
	if ox < 0 {
		ox = r.BoardOffsetX + r.Config.BoardWidth*r.Config.CellWidth + 22
	}
	oy := r.BoardOffsetY + 2

	helpStyle := tcell.StyleDefault.Foreground(tcell.NewRGBColor(120, 120, 150))
	keyStyle := tcell.StyleDefault.Foreground(tcell.NewRGBColor(0, 200, 255)).Bold(true)

	drawString(r.Screen, ox, oy, styleText, "── 操作说明 ──")
	oy += 2
	controls := []struct{ key, desc string }{
		{"←  →", "左右移动"},
		{"  ↓ ", "加速下落"},
		{"  ↑ ", "旋转方块"},
		{"Space", "直接到底"},
		{"  P  ", "暂停游戏"},
		{"  Q  ", "退出游戏"},
	}
	for _, ctrl := range controls {
		drawString(r.Screen, ox, oy, keyStyle, ctrl.key)
		drawString(r.Screen, ox+6, oy, helpStyle, ctrl.desc)
		oy++
	}
}

func (r *Renderer) drawGameOver() {
	ox, oy := r.BoardOffsetX, r.BoardOffsetY
	cx := ox + (r.Config.BoardWidth*r.Config.CellWidth)/2
	cy := oy + (r.Config.BoardHeight*r.Config.CellHeight)/2

	boxW := 22
	boxH := 5
	bx := cx - boxW/2
	by := cy - boxH/2
	boxStyle := tcell.StyleDefault.Background(tcell.NewRGBColor(30, 10, 10)).Foreground(tcell.NewRGBColor(255, 60, 60))
	for row := 0; row < boxH; row++ {
		for c := 0; c < boxW; c++ {
			r.Screen.SetContent(bx+c, by+row, ' ', nil, boxStyle)
		}
	}
	// 稍微调整居中位置，中文字符占用的视觉宽度更大
	drawString(r.Screen, cx-4, cy-1, styleGameOver, "游 戏 结 束")
	drawString(r.Screen, cx-8, cy+1, styleText, "按 R 键重新开始")
}

func (r *Renderer) drawPaused() {
	ox, oy := r.BoardOffsetX, r.BoardOffsetY
	cx := ox + (r.Config.BoardWidth*r.Config.CellWidth)/2
	cy := oy + (r.Config.BoardHeight*r.Config.CellHeight)/2

	boxW := 16
	boxH := 3
	bx := cx - boxW/2
	by := cy - boxH/2
	boxStyle := tcell.StyleDefault.Background(tcell.NewRGBColor(30, 30, 10)).Foreground(tcell.NewRGBColor(255, 200, 0))
	for row := 0; row < boxH; row++ {
		for c := 0; c < boxW; c++ {
			r.Screen.SetContent(bx+c, by+row, ' ', nil, boxStyle)
		}
	}
	drawString(r.Screen, cx-5, cy, stylePause, "⏸ 游戏暂停")
}
