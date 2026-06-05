package game

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"tetris/internal/config"
)

type Cell struct {
	Filled bool
	Color  tcell.Color
}

type Piece struct {
	TetrominoIdx int
	Rotation     int
	X, Y         int
}

type Engine struct {
	Mu           sync.Mutex
	Config       *config.Config
	Board        [][]Cell
	Current      Piece
	Next         int
	Score        int
	Level        int
	Lines        int
	GameOver     bool
	Paused       bool
	DropInterval time.Duration
	GhostY       int
}

func NewEngine(cfg *config.Config) *Engine {
	e := &Engine{
		Config:       cfg,
		Level:        1,
		DropInterval: cfg.InitialSpeed,
	}
	// 动态初始化棋盘大小
	e.Board = make([][]Cell, cfg.BoardHeight)
	for i := range e.Board {
		e.Board[i] = make([]Cell, cfg.BoardWidth)
	}

	e.Next = rand.Intn(len(Tetrominoes))
	e.spawnPiece()
	return e
}

func (e *Engine) spawnPiece() {
	e.Current = Piece{
		TetrominoIdx: e.Next,
		Rotation:     0,
		X:            e.Config.BoardWidth/2 - 2,
		Y:            -1,
	}
	e.Next = rand.Intn(len(Tetrominoes))
	if e.collides(e.Current) {
		e.GameOver = true
	}
	e.UpdateGhost()
}

func (e *Engine) shape(p Piece) [4][4]int {
	return Tetrominoes[p.TetrominoIdx].Rotations[p.Rotation]
}

func (e *Engine) collides(p Piece) bool {
	shape := e.shape(p)
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if shape[r][c] == 0 {
				continue
			}
			nx, ny := p.X+c, p.Y+r
			if nx < 0 || nx >= e.Config.BoardWidth || ny >= e.Config.BoardHeight {
				return true
			}
			if ny >= 0 && e.Board[ny][nx].Filled {
				return true
			}
		}
	}
	return false
}

func (e *Engine) lockPiece() {
	shape := e.shape(e.Current)
	col := Tetrominoes[e.Current.TetrominoIdx].Color
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if shape[r][c] == 0 {
				continue
			}
			nx, ny := e.Current.X+c, e.Current.Y+r
			if ny >= 0 && ny < e.Config.BoardHeight && nx >= 0 && nx < e.Config.BoardWidth {
				e.Board[ny][nx] = Cell{Filled: true, Color: col}
			}
		}
	}
	e.clearLines()
	e.spawnPiece()
}

func (e *Engine) clearLines() {
	cleared := 0
	for r := e.Config.BoardHeight - 1; r >= 0; r-- {
		full := true
		for c := 0; c < e.Config.BoardWidth; c++ {
			if !e.Board[r][c].Filled {
				full = false
				break
			}
		}
		if full {
			cleared++
			for rr := r; rr > 0; rr-- {
				copy(e.Board[rr], e.Board[rr-1])
			}
			e.Board[0] = make([]Cell, e.Config.BoardWidth)
			r++ // 重新检查这一行
		}
	}
	if cleared > 0 {
		points := []int{0, 100, 300, 500, 800}
		if cleared > 4 {
			cleared = 4
		}
		e.Score += points[cleared] * e.Level
		e.Lines += cleared
		e.Level = e.Lines/10 + 1

		// 加速逻辑
		e.DropInterval = time.Duration(800-((e.Level-1)*50)) * time.Millisecond
		if e.DropInterval < 100*time.Millisecond {
			e.DropInterval = 100 * time.Millisecond
		}
	}
}

func (e *Engine) MoveLeft() {
	p := e.Current
	p.X--
	if !e.collides(p) {
		e.Current = p
		e.UpdateGhost()
	}
}

func (e *Engine) MoveRight() {
	p := e.Current
	p.X++
	if !e.collides(p) {
		e.Current = p
		e.UpdateGhost()
	}
}

func (e *Engine) MoveDown() bool {
	p := e.Current
	p.Y++
	if !e.collides(p) {
		e.Current = p
		e.UpdateGhost()
		return true
	}
	return false
}

func (e *Engine) HardDrop() {
	for e.MoveDown() {
		e.Score += 2
	}
	e.lockPiece()
}

func (e *Engine) Rotate() {
	p := e.Current
	p.Rotation = (p.Rotation + 1) % 4
	for _, dx := range []int{0, -1, 1, -2, 2} {
		pp := p
		pp.X += dx
		if !e.collides(pp) {
			e.Current = pp
			e.UpdateGhost()
			return
		}
	}
}

func (e *Engine) UpdateGhost() {
	p := e.Current
	for {
		p.Y++
		if e.collides(p) {
			p.Y--
			break
		}
	}
	e.GhostY = p.Y
}

func (e *Engine) Tick() {
	if e.GameOver || e.Paused {
		return
	}
	if !e.MoveDown() {
		e.lockPiece()
	}
}

func (e *Engine) Reset() {
	*e = *NewEngine(e.Config)
}

func (e *Engine) TogglePause() {
	e.Paused = !e.Paused
}
