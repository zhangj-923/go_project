package board

import "testing"

func TestBoardWinCondition(t *testing.T) {
	b := NewBoard(15)

	// Horizontal Win
	for i := 0; i < 5; i++ {
		b.PlaceStone(i, 0, Black)
	}

	if b.Winner != Black {
		t.Errorf("Expected Black to win, but got %v", b.Winner)
	}

	// Reset
	b = NewBoard(15)

	// Vertical Win
	for i := 0; i < 5; i++ {
		b.PlaceStone(0, i, White)
	}

	if b.Winner != White {
		t.Errorf("Expected White to win, but got %v", b.Winner)
	}
}

func TestBoardUndo(t *testing.T) {
	b := NewBoard(15)
	b.PlaceStone(7, 7, Black)

	if b.Grid[7][7] != Black {
		t.Errorf("Stone not placed correctly")
	}

	b.Undo()
	if b.Grid[7][7] != Empty {
		t.Errorf("Undo failed, expected Empty but got %v", b.Grid[7][7])
	}
}
