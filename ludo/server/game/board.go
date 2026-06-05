package game

func ColorIndex(c Color) int {
	switch c {
	case Red:
		return 0
	case Yellow:
		return 1
	case Blue:
		return 2
	case Green:
		return 3
	}
	return -1
}

func CellColor(pos int) Color {
	if pos >= 0 && pos <= 51 {
		switch pos % 4 {
		case 0:
			return Red
		case 1:
			return Green
		case 2:
			return Yellow
		case 3:
			return Blue
		}
	}
	if pos >= 100 && pos <= 105 {
		return Red
	}
	if pos >= 110 && pos <= 115 {
		return Green
	}
	if pos >= 120 && pos <= 125 {
		return Yellow
	}
	if pos >= 130 && pos <= 135 {
		return Blue
	}
	return ""
}

func IsFlyCell(c Color, pos int) (bool, int) {
	if c == Red && pos == 16 {
		return true, 28
	}
	if c == Green && pos == 29 {
		return true, 41
	}
	if c == Yellow && pos == 42 {
		return true, 2
	}
	if c == Blue && pos == 3 {
		return true, 15
	}
	return false, 0
}

func IsWinPos(c Color, pos int) bool {
	switch c {
	case Red:
		return pos == 105
	case Green:
		return pos == 115
	case Yellow:
		return pos == 125
	case Blue:
		return pos == 135
	}
	return false
}
