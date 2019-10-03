package eit

import (
	"github.com/gdamore/tcell"
	"io/ioutil"
	"os"
	"strings"
)

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

func NewBuffer(config *Config) (buf *Buffer, err error) {
	cursor := new(Cursor)
	buf = &Buffer{
		CurrCursor: cursor,
		Config:     config,
	}
	cursor.Buffer = buf
	err = buf.LoadFromFile(config.FilePath)
	return buf, err
}

func (buffer *Buffer) LoadFromFile(filepath string) error {

	if filepath == "" {
		return nil
	}

	// file not exist
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	left := 0
	for cursor, b := range data {
		if b == byte('\n') {
			buffer.AddLine([]rune(string(data[left:cursor])))
			left = cursor + 1
		}
	}

	if left != len(data) {
		buffer.AddLine([]rune(string(data[left:])))
	}
	return nil
}

func (buffer *Buffer) SaveAs() (err error) {
	if buffer.Config.FilePath == "" {
		return
	}
	err = ioutil.WriteFile(buffer.Config.FilePath, []byte(buffer.String()), 0644)
	return err
}

func (buffer *Buffer) AddLine(line []rune) {
	buffer.Runes = append(buffer.Runes, line)
	buffer.Lines = len(buffer.Runes) - 1
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
		if buffer.CurrCursor.AtEndOfLine() {
			buffer.Runes[currentLine] = buffer.Runes[currentLine][:buffer.CurrCursor.X-1]
		} else {
			buffer.Runes[currentLine] = append(
				buffer.Runes[currentLine][:buffer.CurrCursor.X-1],
				buffer.Runes[currentLine][buffer.CurrCursor.X:]...)
		}

		buffer.CurrCursor.MoveLeft()
	}
}

func (buffer *Buffer) NewLine() {
	buffer.CurrCursor.NewLine()

	buffer.Runes = append(buffer.Runes, []rune{})
	if buffer.Lines == 0 {
		buffer.Lines += 1
		return
	}

	copy(buffer.Runes[buffer.CurrCursor.Y+1:], buffer.Runes[buffer.CurrCursor.Y:])
	buffer.Runes[buffer.CurrCursor.Y] = []rune{}
	buffer.Lines += 1
}

func (buffer *Buffer) Draw(screen tcell.Screen) {
	screen.Clear()
	for y, line := range buffer.Runes {
		for x, char := range line {
			screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}
	}
	screen.ShowCursor(buffer.CurrCursor.X, buffer.CurrCursor.Y)
	screen.Show()
}

func (buffer *Buffer) String() string {
	s := make([]string, buffer.Lines+1)
	for i, rs := range buffer.Runes {
		s[i] = string(rs)
	}

	return strings.Join(s, "\n")
}
