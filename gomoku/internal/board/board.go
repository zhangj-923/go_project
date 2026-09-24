package board

const DefaultSize = 15

type Player int

const (
	Empty Player = iota
	Black        // 先手
	White        // 后手
)

func (p Player) String() string {
	switch p {
	case Black:
		return "Black"
	case White:
		return "White"
	default:
		return "Empty"
	}
}

type Board struct {
	Size         int
	Grid         [][]Player
	LastMove     Move
	MoveHistory  []Move
	Winner       Player
	WinningLines []Point
}

type Point struct {
	X, Y int
}

type Move struct {
	P      Point
	Player Player
}

func NewBoard(size int) *Board {
	if size <= 0 {
		size = DefaultSize
	}
	grid := make([][]Player, size)
	for i := range grid {
		grid[i] = make([]Player, size)
	}
	return &Board{
		Size:        size,
		Grid:        grid,
		LastMove:    Move{Point{-1, -1}, Empty},
		MoveHistory: make([]Move, 0),
	}
}

func (b *Board) IsValidMove(x, y int) bool {
	if x < 0 || x >= b.Size || y < 0 || y >= b.Size {
		return false
	}
	return b.Grid[y][x] == Empty
}

func (b *Board) PlaceStone(x, y int, p Player) bool {
	if !b.IsValidMove(x, y) {
		return false
	}
	b.Grid[y][x] = p
	b.LastMove = Move{Point{x, y}, p}
	b.MoveHistory = append(b.MoveHistory, b.LastMove)
	b.checkWin(x, y, p)
	return true
}

func (b *Board) Undo() {
	if len(b.MoveHistory) == 0 {
		return
	}
	last := b.MoveHistory[len(b.MoveHistory)-1]
	b.MoveHistory = b.MoveHistory[:len(b.MoveHistory)-1]
	b.Grid[last.P.Y][last.P.X] = Empty
	if len(b.MoveHistory) > 0 {
		b.LastMove = b.MoveHistory[len(b.MoveHistory)-1]
	} else {
		b.LastMove = Move{Point{-1, -1}, Empty}
	}
	b.Winner = Empty
	b.WinningLines = nil
}

func (b *Board) checkWin(x, y int, p Player) {
	directions := []Point{
		{1, 0},  // Horizontal
		{0, 1},  // Vertical
		{1, 1},  // Diagonal \
		{1, -1}, // Diagonal /
	}

	for _, d := range directions {
		count := 1
		line := []Point{{x, y}}

		// Forward
		for i := 1; i < 5; i++ {
			nx, ny := x+d.X*i, y+d.Y*i
			if nx >= 0 && nx < b.Size && ny >= 0 && ny < b.Size && b.Grid[ny][nx] == p {
				count++
				line = append(line, Point{nx, ny})
			} else {
				break
			}
		}

		// Backward
		for i := 1; i < 5; i++ {
			nx, ny := x-d.X*i, y-d.Y*i
			if nx >= 0 && nx < b.Size && ny >= 0 && ny < b.Size && b.Grid[ny][nx] == p {
				count++
				line = append(line, Point{nx, ny})
			} else {
				break
			}
		}

		if count >= 5 {
			b.Winner = p
			b.WinningLines = line
			return
		}
	}
}

func (b *Board) IsFull() bool {
	for y := 0; y < b.Size; y++ {
		for x := 0; x < b.Size; x++ {
			if b.Grid[y][x] == Empty {
				return false
			}
		}
	}
	return true
}
