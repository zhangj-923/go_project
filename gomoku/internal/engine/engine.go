package engine

import (
	"image/color"
	"math/rand"
	"time"

	"gomoku/internal/ai"
	"gomoku/internal/board"
	"gomoku/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

const (
	Menu GameState = iota
	Playing
	GameOver
)

type Engine struct {
	state       GameState
	board       *board.Board
	aiPlayer    *ai.AI
	humanPlayer board.Player
	aiSide      board.Player
	turn        board.Player
	uiState     *ui.UIState

	screenWidth  int
	screenHeight int

	hoverX, hoverY int
}

func NewEngine() *Engine {
	rand.Seed(time.Now().UnixNano())
	ui.InitFonts()
	return &Engine{
		state:       Menu,
		board:       board.NewBoard(),
		humanPlayer: board.Black,
		aiSide:      board.White,
		turn:        board.Black,
		uiState:     ui.NewUIState(),
		aiPlayer:    ai.NewAI(ai.Normal),
	}
}

func (e *Engine) Update() error {
	e.screenWidth, e.screenHeight = ebiten.WindowSize()

	e.uiState.Update(e.screenWidth, e.screenHeight)

	switch e.state {
	case Menu:
		e.updateMenu()
	case Playing:
		e.updatePlaying()
	case GameOver:
		e.updateGameOver()
	}

	return nil
}

func (e *Engine) updateMenu() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		// Start game button
		cx, cy := e.screenWidth/2, e.screenHeight/2
		if mx >= cx-120 && mx <= cx+120 && my >= cy-25 && my <= cy+25 {
			e.startNewGame()
		}

		// Difficulty buttons
		if my >= cy+50 && my <= cy+90 {
			if mx >= cx-180 && mx < cx-70 {
				e.aiPlayer.Difficulty = ai.Easy
			} else if mx >= cx-60 && mx < cx+50 {
				e.aiPlayer.Difficulty = ai.Normal
			} else if mx >= cx+60 && mx <= cx+170 {
				e.aiPlayer.Difficulty = ai.Hard
			}
		}
	}
}

func (e *Engine) startNewGame() {
	e.board = board.NewBoard()
	e.turn = board.Black
	e.state = Playing
}

func (e *Engine) updatePlaying() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		e.state = Menu
		return
	}

	if e.turn == e.aiSide {
		// AI Turn
		move := e.aiPlayer.GetBestMove(e.board, e.aiSide)
		if move.X != -1 && move.Y != -1 {
			e.board.PlaceStone(move.X, move.Y, e.aiSide)
		}
		e.checkGameOver()
		e.turn = e.humanPlayer
	} else {
		// Human Turn
		mx, my := ebiten.CursorPosition()
		bx, by := e.uiState.ScreenToBoard(mx, my, e.screenWidth, e.screenHeight)

		e.hoverX, e.hoverY = bx, by

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if e.board.IsValidMove(bx, by) {
				e.board.PlaceStone(bx, by, e.humanPlayer)
				e.checkGameOver()
				e.turn = e.aiSide
			}
		}
	}
}

func (e *Engine) updateGameOver() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		e.state = Menu
	}
}

func (e *Engine) checkGameOver() {
	if e.board.Winner != board.Empty || e.board.IsFull() {
		e.state = GameOver
	}
}

func (e *Engine) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{240, 230, 210, 255}) // Background

	switch e.state {
	case Menu:
		e.uiState.DrawMenu(screen, e.screenWidth, e.screenHeight, e.aiPlayer.Difficulty)
	case Playing:
		e.uiState.DrawBoard(screen, e.board, e.screenWidth, e.screenHeight, e.hoverX, e.hoverY, e.turn == e.humanPlayer)
		e.uiState.DrawHUD(screen, e.turn, e.aiPlayer.Difficulty)
	case GameOver:
		e.uiState.DrawBoard(screen, e.board, e.screenWidth, e.screenHeight, -1, -1, false)
		e.uiState.DrawGameOver(screen, e.board.Winner, e.screenWidth, e.screenHeight)
	}
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
