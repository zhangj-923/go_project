package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type RoomState string

const (
	StateWaiting  RoomState = "waiting"
	StatePlaying  RoomState = "playing"
	StateGameOver RoomState = "game_over"
)

type Client interface {
	Send(msg []byte)
	ID() string
	Name() string
}

type Room struct {
	ID         string
	GameMode   string
	MaxPlayers int
	AICount    int
	State      RoomState

	Board     *BoardState
	Clients   map[string]Client
	TurnIndex int

	CurrentDice   int
	ContinuousSix int
	HasRolled     bool

	register   chan Client
	unregister chan Client
	actionChan chan ActionMessage

	mu sync.Mutex
}

type ActionMessage struct {
	Client Client
	Msg    Message
}

func NewRoom(id, mode string, maxP, aiCount int) *Room {
	return &Room{
		ID:         id,
		GameMode:   mode,
		MaxPlayers: maxP,
		AICount:    aiCount,
		State:      StateWaiting,
		Board:      &BoardState{Players: make([]Player, 0)},
		Clients:    make(map[string]Client),
		register:   make(chan Client),
		unregister: make(chan Client),
		actionChan: make(chan ActionMessage),
	}
}

func (r *Room) RegisterChan() chan<- Client      { return r.register }
func (r *Room) UnregisterChan() chan<- Client    { return r.unregister }
func (r *Room) ActionChan() chan<- ActionMessage { return r.actionChan }

func (r *Room) Run() {
	fmt.Println("Room Run loop started for", r.ID)
	for {
		select {
		case client := <-r.register:
			fmt.Println("Room received client to register:", client.ID())
			r.handleJoin(client)
		case client := <-r.unregister:
			r.handleLeave(client)
		case action := <-r.actionChan:
			r.handleAction(action)
		}
	}
}

func (r *Room) BroadcastMsg(msgType MsgType, payload interface{}) {
	msg := Message{Type: msgType, Payload: payload}
	b, _ := json.Marshal(msg)
	for _, client := range r.Clients {
		client.Send(b)
	}
}

func (r *Room) handleJoin(client Client) {
	if r.State != StateWaiting {
		// Cannot join
		return
	}
	r.Clients[client.ID()] = client

	color := Colors[len(r.Board.Players)]
	player := NewPlayer(client.ID(), client.Name(), color, false)
	r.Board.Players = append(r.Board.Players, *player)

	r.BroadcastMsg(MsgRoomJoined, PayloadRoomJoined{
		RoomID:   r.ID,
		Players:  r.Board.Players, // wait, need pointers? No, copy is fine for broadcast
		MyColor:  color,
		MyID:     client.ID(),
		GameMode: r.GameMode,
	})

	// Add AIs if needed when first player joins or auto fill
	// Simplification: Wait for owner to click start.
}

func (r *Room) handleLeave(client Client) {
	delete(r.Clients, client.ID())
	// Set player connected = false
	for i := range r.Board.Players {
		if r.Board.Players[i].ID == client.ID() {
			r.Board.Players[i].Connected = false
		}
	}
	r.BroadcastMsg(MsgPlayerLeft, client.ID())
}

func (r *Room) handleAction(action ActionMessage) {
	switch action.Msg.Type {
	case MsgStartGame:
		r.startGame()
	case MsgRollDice:
		r.handleRollDice(action.Client)
	case MsgMovePiece:
		var payload PayloadMovePiece
		b, _ := json.Marshal(action.Msg.Payload)
		json.Unmarshal(b, &payload)
		r.handleMovePiece(action.Client, payload.PieceID)
	case "internal_next_turn":
		r.nextTurn()
	case "internal_ai_move":
		var payload PayloadMovePiece
		b, _ := json.Marshal(action.Msg.Payload)
		json.Unmarshal(b, &payload)
		r.handleMovePiece(action.Client, payload.PieceID)
	}
}

func (r *Room) startGame() {
	if r.State != StateWaiting {
		return
	}
	// Add AI players
	for i := 0; i < r.AICount; i++ {
		color := Colors[len(r.Board.Players)]
		ai := NewPlayer(fmt.Sprintf("AI_%d", i), fmt.Sprintf("电脑 %d", i+1), color, true)
		r.Board.Players = append(r.Board.Players, *ai)
	}

	r.State = StatePlaying
	r.TurnIndex = 0
	r.CurrentDice = 0
	r.ContinuousSix = 0
	r.HasRolled = false

	r.BroadcastMsg(MsgGameStarted, PayloadGameStarted{
		BoardState:  r.Board,
		CurrentTurn: r.Board.Players[r.TurnIndex].Color,
	})

	r.checkAITurn()
}

func (r *Room) handleRollDice(client Client) {
	if r.State != StatePlaying {
		return
	}
	currPlayer := &r.Board.Players[r.TurnIndex]
	if currPlayer.ID != client.ID() || r.HasRolled {
		return
	}

	r.CurrentDice = RollDice()
	r.HasRolled = true

	if r.CurrentDice == 6 {
		r.ContinuousSix++
	} else {
		r.ContinuousSix = 0
	}

	// 连六检查 (连续三次6返回机场)
	if r.ContinuousSix == 3 {
		// 返回机场逻辑 - 简单处理：跳过回合
		r.ContinuousSix = 0
		r.BroadcastMsg(MsgDiceRolled, PayloadDiceRolled{
			PlayerID:      currPlayer.ID,
			Color:         currPlayer.Color,
			Value:         6,
			MovablePieces: []int{},
			ContinuousSix: 3,
		})
		r.nextTurn()
		return
	}

	movable := GetMovablePieces(r.Board, currPlayer.Color, r.CurrentDice)

	r.BroadcastMsg(MsgDiceRolled, PayloadDiceRolled{
		PlayerID:      currPlayer.ID,
		Color:         currPlayer.Color,
		Value:         r.CurrentDice,
		MovablePieces: movable,
		ContinuousSix: r.ContinuousSix,
	})

	if len(movable) == 0 {
		// No moves possible, delay and next turn
		go func() {
			time.Sleep(1 * time.Second)
			r.actionChan <- ActionMessage{Msg: Message{Type: "internal_next_turn"}}
		}()
	} else if currPlayer.IsAI {
		// AI move
		go func(playerID string, movables []int) {
			time.Sleep(1 * time.Second)
			pieceID := movables[rand.Intn(len(movables))]
			r.actionChan <- ActionMessage{
				Client: &DummyClient{id: playerID},
				Msg:    Message{Type: "internal_ai_move", Payload: map[string]interface{}{"pieceId": float64(pieceID)}},
			}
		}(currPlayer.ID, movable)
	}
}

func (r *Room) handleMovePiece(client Client, pieceID int) {
	if r.State != StatePlaying || !r.HasRolled {
		return
	}
	currPlayer := &r.Board.Players[r.TurnIndex]
	if currPlayer.ID != client.ID() {
		return
	}

	// Move the piece
	var piece *Piece
	for i := range currPlayer.Pieces {
		if currPlayer.Pieces[i].ID == pieceID {
			piece = &currPlayer.Pieces[i]
			break
		}
	}
	if piece == nil {
		return
	}

	oldPos := piece.Position
	newPos, events := CalculateMove(r.Board, piece, r.CurrentDice)
	if newPos == -1 { // Invalid move
		return
	}

	// Apply events (Eat)
	for _, e := range events {
		if e.Type == EventEat {
			// Find eaten piece and send back to base
			// We didn't store enemy color in event, let's just search board at e.Position (which is -1 for eaten piece)
			// Wait, CalculateMove sets EventEat.Position to -1.
			// Let's modify CalculateMove to pass color. Or just re-scan here.
			r.eatPiecesAt(oldPos, newPos, currPlayer.Color)
		}
	}

	piece.Position = newPos

	r.BroadcastMsg(MsgPieceMoved, PayloadPieceMoved{
		PieceID: pieceID,
		Color:   currPlayer.Color,
		From:    oldPos,
		To:      newPos,
		Events:  events,
	})

	// Check win
	if r.checkWin(currPlayer) {
		r.State = StateGameOver
		r.BroadcastMsg(MsgGameOver, PayloadGameOver{
			WinnerID: currPlayer.ID,
			Color:    currPlayer.Color,
		})
		return
	}

	// Extra roll for 6
	if r.CurrentDice == 6 {
		r.HasRolled = false
		r.BroadcastMsg(MsgTurnChanged, PayloadTurnChanged{
			PlayerID:  currPlayer.ID,
			Color:     currPlayer.Color,
			ExtraRoll: true,
		})
		r.checkAITurn()
	} else {
		r.nextTurn()
	}
}

func (r *Room) eatPiecesAt(oldPos, newPos int, myColor Color) {
	for i := range r.Board.Players {
		p := &r.Board.Players[i]
		if p.Color == myColor {
			continue
		}
		for j := range p.Pieces {
			if p.Pieces[j].Position == newPos && newPos >= 0 && newPos <= 51 {
				p.Pieces[j].Position = -1 // Back to base
			}
		}
	}
}

func (r *Room) nextTurn() {
	r.TurnIndex = (r.TurnIndex + 1) % len(r.Board.Players)
	r.HasRolled = false
	r.CurrentDice = 0

	// Reset continuous six for next player if it's not the same player
	// Wait, we only advance turn if they didn't roll 6. So continuous six should be 0.
	r.ContinuousSix = 0

	currPlayer := &r.Board.Players[r.TurnIndex]
	r.BroadcastMsg(MsgTurnChanged, PayloadTurnChanged{
		PlayerID:  currPlayer.ID,
		Color:     currPlayer.Color,
		ExtraRoll: false,
	})

	r.checkAITurn()
}

func (r *Room) checkAITurn() {
	currPlayer := r.Board.Players[r.TurnIndex]
	if currPlayer.IsAI {
		go func(player Player) {
			time.Sleep(1 * time.Second)
			r.actionChan <- ActionMessage{Client: &DummyClient{id: player.ID}, Msg: Message{Type: MsgRollDice}}
		}(currPlayer)
	}
}

func (r *Room) checkWin(player *Player) bool {
	winCount := 0
	for _, p := range player.Pieces {
		if IsWinPos(player.Color, p.Position) {
			winCount++
		}
	}
	return winCount == 4
}

// DummyClient for AI
type DummyClient struct {
	id string
}

func (c *DummyClient) Send(msg []byte) {}
func (c *DummyClient) ID() string      { return c.id }
func (c *DummyClient) Name() string    { return c.id }
