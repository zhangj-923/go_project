package engine

import (
	"image/color"
	"math/rand"
	"time"

	"gomoku/internal/ai"
	"gomoku/internal/board"
	"gomoku/internal/i18n"
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

type GameMode int

const (
	ModeVsAI GameMode = iota
	ModeLocalPvP
	ModeOnline
)

type Engine struct {
	state       GameState
	mode        GameMode
	boardSize   int
	lang        i18n.Language
	board       *board.Board
	aiPlayer    *ai.AI
	humanPlayer board.Player
	aiSide      board.Player
	turn        board.Player
	uiState     *ui.UIState

	screenWidth  int
	screenHeight int

	hoverX, hoverY int

	onlineRole   string
	onlineRoomID string

	touchIDs []ebiten.TouchID
}

func NewEngine() *Engine {
	rand.Seed(time.Now().UnixNano())
	ui.InitFonts()
	defaultSize := 15
	e := &Engine{
		state:        Menu,
		mode:         ModeVsAI,
		boardSize:    defaultSize,
		lang:         i18n.ZH_CN,
		board:        board.NewBoard(defaultSize),
		humanPlayer:  board.Black,
		aiSide:       board.White,
		turn:         board.Black,
		uiState:      ui.NewUIState(),
		aiPlayer:     ai.NewAI(ai.Normal),
		screenWidth:  1280,
		screenHeight: 900,
		touchIDs:     make([]ebiten.TouchID, 0, 8),
	}
	e.initJSBridge()
	return e
}

func (e *Engine) isClicked() (bool, int, int) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		return true, mx, my
	}
	e.touchIDs = inpututil.AppendJustPressedTouchIDs(e.touchIDs[:0])
	if len(e.touchIDs) > 0 {
		tx, ty := ebiten.TouchPosition(e.touchIDs[0])
		return true, tx, ty
	}
	return false, 0, 0
}

func (e *Engine) Update() error {
	currentSize := e.boardSize
	if e.board != nil {
		currentSize = e.board.Size
	}
	e.uiState.Update(e.screenWidth, e.screenHeight, currentSize, e.mode == ModeOnline)

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
	if clicked, mx, my := e.isClicked(); clicked {
		cx, cy := float32(e.screenWidth)/2, float32(e.screenHeight)/2
		fmx, fmy := float32(mx), float32(my)

		// 1. Language Toggle Button (Top-Right)
		langBtnW := float32(130)
		langBtnH := float32(36)
		langBtnX := float32(e.screenWidth) - langBtnW - 20
		langBtnY := float32(20)
		if fmx >= langBtnX && fmx <= langBtnX+langBtnW && fmy >= langBtnY && fmy <= langBtnY+langBtnH {
			e.lang = i18n.NextLanguage(e.lang)
			return
		}

		// 2. Start Game Button
		startW := float32(260)
		startH := float32(52)
		startX := cx - startW/2
		startY := cy - 110
		if fmx >= startX && fmx <= startX+startW && fmy >= startY && fmy <= startY+startH {
			e.startNewGame()
			return
		}

		// 3. Mode Buttons
		modeY := cy - 35
		if fmy >= modeY && fmy <= modeY+36 {
			// Mode 0: Vs AI
			if fmx >= cx-115 && fmx <= cx+25 {
				e.mode = ModeVsAI
			}
			// Mode 1: PvP
			if fmx >= cx+40 && fmx <= cx+180 {
				e.mode = ModeLocalPvP
			}
		}

		// 4. Board Size Buttons
		sizeY := cy + 20
		if fmy >= sizeY && fmy <= sizeY+36 {
			if fmx >= cx-95 && fmx <= cx {
				e.boardSize = 11
			} else if fmx >= cx+10 && fmx <= cx+105 {
				e.boardSize = 15
			} else if fmx >= cx+115 && fmx <= cx+210 {
				e.boardSize = 19
			}
		}

		// 5. Difficulty Buttons (Only if Vs AI)
		if e.mode == ModeVsAI {
			diffY := cy + 75
			if fmy >= diffY && fmy <= diffY+36 {
				if fmx >= cx-95 && fmx < cx-25 {
					e.aiPlayer.Difficulty = ai.Easy
				} else if fmx >= cx-17 && fmx < cx+53 {
					e.aiPlayer.Difficulty = ai.Normal
				} else if fmx >= cx+61 && fmx < cx+131 {
					e.aiPlayer.Difficulty = ai.Hard
				} else if fmx >= cx+139 && fmx <= cx+209 {
					e.aiPlayer.Difficulty = ai.Hell
				}
			}
		}
	}
}

func (e *Engine) startNewGame() {
	e.board = board.NewBoard(e.boardSize)
	e.turn = board.Black
	e.state = Playing
}

func (e *Engine) updatePlaying() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		e.state = Menu
		if e.mode == ModeOnline {
			e.notifyReturnMenuToJS()
		}
		return
	}

	if e.mode == ModeOnline {
		mx, my := ebiten.CursorPosition()
		bx, by := e.uiState.ScreenToBoard(mx, my, e.screenWidth, e.screenHeight, e.board.Size)
		e.hoverX, e.hoverY = bx, by

		isMyTurn := (e.onlineRole == "black" && e.turn == board.Black) || (e.onlineRole == "white" && e.turn == board.White)
		if isMyTurn {
			if clicked, cx, cy := e.isClicked(); clicked {
				cbx, cby := e.uiState.ScreenToBoard(cx, cy, e.screenWidth, e.screenHeight, e.board.Size)
				if e.board.IsValidMove(cbx, cby) {
					e.sendLocalMoveToJS(cbx, cby)
				}
			}
		}
		return
	}

	if e.mode == ModeVsAI {
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
			bx, by := e.uiState.ScreenToBoard(mx, my, e.screenWidth, e.screenHeight, e.board.Size)
			e.hoverX, e.hoverY = bx, by

			if clicked, cx, cy := e.isClicked(); clicked {
				cbx, cby := e.uiState.ScreenToBoard(cx, cy, e.screenWidth, e.screenHeight, e.board.Size)
				if e.board.IsValidMove(cbx, cby) {
					e.board.PlaceStone(cbx, cby, e.humanPlayer)
					e.checkGameOver()
					e.turn = e.aiSide
				}
			}
		}
	} else {
		// Local 2-Player (PvP)
		mx, my := ebiten.CursorPosition()
		bx, by := e.uiState.ScreenToBoard(mx, my, e.screenWidth, e.screenHeight, e.board.Size)
		e.hoverX, e.hoverY = bx, by

		if clicked, cx, cy := e.isClicked(); clicked {
			cbx, cby := e.uiState.ScreenToBoard(cx, cy, e.screenWidth, e.screenHeight, e.board.Size)
			if e.board.IsValidMove(cbx, cby) {
				e.board.PlaceStone(cbx, cby, e.turn)
				e.checkGameOver()
				if e.turn == board.Black {
					e.turn = board.White
				} else {
					e.turn = board.Black
				}
			}
		}
	}
}

func (e *Engine) updateGameOver() {
	if e.mode == ModeOnline {
		// 在联机对战模式下，对局结束时不因随意点击棋盘而误退回本地菜单
		// 玩家应通过重赛、重开或主动退出按钮进行操作
		return
	}
	if clicked, _, _ := e.isClicked(); clicked {
		e.state = Menu
	}
}

func (e *Engine) checkGameOver() {
	if e.board.Winner != board.Empty || e.board.IsFull() {
		e.state = GameOver
	}
}

func (e *Engine) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{242, 235, 220, 255}) // Warm background

	if b := screen.Bounds(); b.Dx() > 0 && b.Dy() > 0 {
		e.screenWidth = b.Dx()
		e.screenHeight = b.Dy()
	}

	switch e.state {
	case Menu:
		e.uiState.DrawMenu(screen, e.screenWidth, e.screenHeight, e.mode == ModeVsAI, e.aiPlayer.Difficulty, e.boardSize, e.lang)
	case Playing:
		isHumanTurn := false
		if e.mode == ModeOnline {
			isHumanTurn = (e.onlineRole == "black" && e.turn == board.Black) || (e.onlineRole == "white" && e.turn == board.White)
		} else if e.mode == ModeVsAI {
			isHumanTurn = e.turn == e.humanPlayer
		} else {
			isHumanTurn = true
		}
		e.uiState.DrawBoard(screen, e.board, e.screenWidth, e.screenHeight, e.hoverX, e.hoverY, isHumanTurn)
		e.uiState.DrawHUD(screen, e.mode == ModeVsAI, e.turn, e.aiPlayer.Difficulty, e.board.Size, e.lang)
	case GameOver:
		e.uiState.DrawBoard(screen, e.board, e.screenWidth, e.screenHeight, -1, -1, false)
		e.uiState.DrawGameOver(screen, e.board.Winner, e.screenWidth, e.screenHeight, e.mode == ModeOnline, e.lang)
	}
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth > 0 && outsideHeight > 0 {
		e.screenWidth = outsideWidth
		e.screenHeight = outsideHeight
	}
	return e.screenWidth, e.screenHeight
}
