package eit

// content of screen
type Buffer struct {
	// current cursor
	CurrCursor *Cursor

	// []rune just equivalent of one line, and a rune is a rune.
	Runes [][]rune

	// how many lines of buffer, start with 0
	Lines int

	Config *Config
}

func NewBuffer(config *Config) *Buffer {
	cursor := new(Cursor)
	buf := &Buffer{
		CurrCursor: cursor,
		Config:     config,
	}
	cursor.Buffer = buf
	return buf
}

// Is current buffer empty
func (buffer *Buffer) IsEmpty() bool {
	return buffer.Lines == 0 && len(buffer.Runes) == 0
}

// Insert a rune behind current cursor
func (buffer *Buffer) InsertRune(r rune) {
	defer buffer.CurrCursor.MoveRight()

	if buffer.IsEmpty() {
		buffer.Runes = append(buffer.Runes, []rune{r})
	} else {
		s := append(buffer.Runes[buffer.CurrCursor.Y], rune(0))
		copy(s[buffer.CurrCursor.X+1:], s[buffer.CurrCursor.X:])
		s[buffer.CurrCursor.X] = r
		buffer.Runes[buffer.CurrCursor.Y] = s
	}
}

func (buffer *Buffer) DeleteLine() {
	if buffer.Lines == 0 {
		buffer.Runes = [][]rune{}
		return
	}

	defer buffer.CurrCursor.MovePrevLine()
	buffer.Runes = append(
		buffer.Runes[:buffer.CurrCursor.Y],
		buffer.Runes[buffer.CurrCursor.Y+1:]...)
	buffer.Lines -= 1
}

func (buffer *Buffer) Delete() {
	if buffer.IsEmpty() {
		return
	}

	currentLine := buffer.CurrCursor.Y

	// 行首删除，跳到上一行，并把下一行内容追加到上一行
	if buffer.CurrCursor.AtBeginOfLine() {
		deletedLine := buffer.Runes[buffer.CurrCursor.Y]
		buffer.DeleteLine()
		if len(deletedLine) == 0 {
			return
		}

		// 追加被删除的行到上一行末尾
		buffer.Runes[buffer.CurrCursor.Y] = append(
			buffer.Runes[buffer.CurrCursor.Y],
			deletedLine...)
	} else {
		// 不在行首直接删除 rune
		buffer.Runes[currentLine] = append(
			buffer.Runes[currentLine][:buffer.CurrCursor.X],
			buffer.Runes[currentLine][buffer.CurrCursor.X+1:]...)

		buffer.CurrCursor.MoveLeft()
	}
}

func (buffer *Buffer) NewLine() {
	defer buffer.CurrCursor.NewLine()

	buffer.Runes = append(buffer.Runes, []rune{})
	buffer.Lines += 1
	if buffer.Lines == 2 {
		return
	}

	copy(buffer.Runes[buffer.CurrCursor.Y+1:], buffer.Runes[buffer.CurrCursor.Y:])
	buffer.Runes[buffer.CurrCursor.Y] = []rune{}
}
