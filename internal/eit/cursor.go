package eit

// The cursor position
type Cursor struct {
	X, Y   int
	Buffer *Buffer
}

// cursor at (0, 0) point
func (cursor *Cursor) IsOriginPoint() bool {
	return cursor.X == 0 && cursor.Y == 0
}

// at line end
func (cursor *Cursor) AtEndOfLine() bool {
	if cursor.Buffer.IsEmpty() {
		return true
	}

	return cursor.X == len(cursor.Buffer.Runes[cursor.Buffer.Lines])
}

func (cursor *Cursor) AtBeginOfLine() bool {
	return cursor.X == 0
}

// at the last line of buffer
func (cursor *Cursor) AtLastLineOfBuffer() bool {
	return cursor.Y == cursor.Buffer.Lines
}

func (cursor *Cursor) AtFirstLineOfBuffer() bool {
	return cursor.Y == 0
}

// cursor at buffer last point
func (cursor *Cursor) IsLastPoint() bool {
	return cursor.AtLastLineOfBuffer() && cursor.AtEndOfLine()
}

// Move cursor left
func (cursor *Cursor) MoveLeft() {
	if cursor.Buffer.IsEmpty() || cursor.IsOriginPoint() {
		return
	}

	if cursor.X == 0 {
		cursor.Y -= 1
		cursor.X = len(cursor.Buffer.Runes[cursor.Y])
	} else {
		cursor.X -= 1
	}
}

func (cursor *Cursor) MoveRight() {
	if cursor.Buffer.IsEmpty() || cursor.IsLastPoint() {
		return
	}

	if cursor.AtEndOfLine() {
		cursor.NewLine()
		return
	}

	cursor.X += 1
}

func (cursor *Cursor) MoveTop() {
	if cursor.AtFirstLineOfBuffer() {
		return
	} else {
		cursor.Y -= 1
	}
}

func (cursor *Cursor) MoveDown() {
	if cursor.AtLastLineOfBuffer() {
		return
	} else {
		cursor.Y += 1
	}
}

func (cursor *Cursor) MovePrevLine() {
	if cursor.AtFirstLineOfBuffer() {
		return
	}

	cursor.Y -= 1
	cursor.X = len(cursor.Buffer.Runes[cursor.Y])
}

func (cursor *Cursor) NewLine() {
	cursor.Y += 1
	cursor.X = 0
}
