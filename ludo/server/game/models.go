package game

type Color string

const (
	Red    Color = "red"
	Yellow Color = "yellow"
	Blue   Color = "blue"
	Green  Color = "green"
)

var Colors = []Color{Red, Yellow, Blue, Green}

type MoveEventType string

const (
	EventMove MoveEventType = "move"
	EventJump MoveEventType = "jump"
	EventFly  MoveEventType = "fly"
	EventEat  MoveEventType = "eat"
	EventHome MoveEventType = "home"
)

type MoveEvent struct {
	Type     MoveEventType `json:"type"`
	Position int           `json:"position"`         // the intermediate or final position
	Target   int           `json:"target,omitempty"` // for eat, which piece was eaten
}

type Piece struct {
	ID       int   `json:"id"` // 0-3 for each player
	Color    Color `json:"color"`
	Position int   `json:"position"` // -1: Base. 0-51: Track. 100+: Home stretch.
}

type Player struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Color     Color   `json:"color"`
	IsAI      bool    `json:"isAi"`
	Pieces    []Piece `json:"pieces"`
	Rank      int     `json:"rank"` // 1, 2, 3, 4. 0 means not finished.
	Connected bool    `json:"connected"`
}

func NewPlayer(id, name string, color Color, isAI bool) *Player {
	p := &Player{
		ID:        id,
		Name:      name,
		Color:     color,
		IsAI:      isAI,
		Pieces:    make([]Piece, 4),
		Connected: true,
	}
	for i := 0; i < 4; i++ {
		p.Pieces[i] = Piece{
			ID:       i,
			Color:    color,
			Position: -1, // -1 means in Hangar/Base
		}
	}
	return p
}

// BoardState for broadcasting
type BoardState struct {
	Players []Player `json:"players"`
}

// Map from Color to starting position on the main track (0-51)
var StartPositions = map[Color]int{
	Red:    0,
	Green:  13,
	Yellow: 26,
	Blue:   39,
}

// Map from Color to the cell BEFORE entering the home stretch
var TurnHomePositions = map[Color]int{
	Red:    50, // Turns to 100 after 50
	Green:  11, // Turns to 110 after 11
	Yellow: 24, // Turns to 120 after 24
	Blue:   37, // Turns to 130 after 37
}

// Base offsets for home stretch
var HomeStretchBases = map[Color]int{
	Red:    100, // 100 - 105 (105 is win)
	Green:  110,
	Yellow: 120,
	Blue:   130,
}

// Fly start cells and target cells
var FlyPaths = map[Color]struct {
	Start int
	End   int
}{
	Red:    {Start: 18, End: 30}, // Red flies from 18 to 30. Wait, standard rules differ.
	Yellow: {Start: 31, End: 43},
	Blue:   {Start: 44, End: 4},
	Green:  {Start: 5, End: 17},
}

// I will refine FlyPaths in board.go
