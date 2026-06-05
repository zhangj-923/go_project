package game

type MsgType string

const (
	// Client -> Server
	MsgCreateRoom MsgType = "create_room"
	MsgJoinRoom   MsgType = "join_room"
	MsgStartGame  MsgType = "start_game"
	MsgRollDice   MsgType = "roll_dice"
	MsgMovePiece  MsgType = "move_piece"
	MsgLeaveRoom  MsgType = "leave_room"

	// Server -> Client
	MsgRoomJoined   MsgType = "room_joined"
	MsgGameStarted  MsgType = "game_started"
	MsgDiceRolled   MsgType = "dice_rolled"
	MsgPieceMoved   MsgType = "piece_moved"
	MsgTurnChanged  MsgType = "turn_changed"
	MsgPlayerJoined MsgType = "player_joined"
	MsgPlayerLeft   MsgType = "player_left"
	MsgGameOver     MsgType = "game_over"
	MsgError        MsgType = "error"
)

type Message struct {
	Type    MsgType     `json:"type"`
	Payload interface{} `json:"payload"`
}

type PayloadCreateRoom struct {
	PlayerName string `json:"playerName"`
	GameMode   string `json:"gameMode"` // "single", "local", "lan"
	MaxPlayers int    `json:"maxPlayers"`
	AICount    int    `json:"aiCount"`
}

type PayloadJoinRoom struct {
	RoomID     string `json:"roomId"`
	PlayerName string `json:"playerName"`
}

type PayloadMovePiece struct {
	PieceID int `json:"pieceId"`
}

type PayloadRoomJoined struct {
	RoomID   string   `json:"roomId"`
	Players  []Player `json:"players"`
	MyColor  Color    `json:"myColor"`
	MyID     string   `json:"myId"`
	GameMode string   `json:"gameMode"`
}

type PayloadGameStarted struct {
	BoardState  *BoardState `json:"boardState"`
	CurrentTurn Color       `json:"currentTurn"`
}

type PayloadDiceRolled struct {
	PlayerID      string `json:"playerId"`
	Color         Color  `json:"color"`
	Value         int    `json:"value"`
	MovablePieces []int  `json:"movablePieces"`
	ContinuousSix int    `json:"continuousSix"`
}

type PayloadPieceMoved struct {
	PieceID int         `json:"pieceId"`
	Color   Color       `json:"color"`
	From    int         `json:"from"`
	To      int         `json:"to"`
	Events  []MoveEvent `json:"events"` // Jump, Fly, Eat
}

type PayloadTurnChanged struct {
	PlayerID  string `json:"playerId"`
	Color     Color  `json:"color"`
	ExtraRoll bool   `json:"extraRoll"`
}

type PayloadGameOver struct {
	WinnerID string   `json:"winnerId"`
	Color    Color    `json:"color"`
	Rankings []Player `json:"rankings"`
}

type PayloadError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
