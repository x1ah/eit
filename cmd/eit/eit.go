package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"os"
)

type Content struct {
	Runes [][]rune
	Lines int32
}

func (text *Content) IsEmpty() bool {
	return text.Lines == 0 && len(text.Runes) == 0
}

func (text *Content) AddRune(char rune) {
	if text.IsEmpty() {
		text.Runes = append(text.Runes, []rune{char})
	} else {
		text.Runes[text.Lines] = append(text.Runes[text.Lines], char)
	}
}

func (text *Content) DeleteRune() {
	if text.IsEmpty() {
		return
	}

	if len(text.Runes[text.Lines]) == 0 {
		text.Runes = text.Runes[:text.Lines]
		text.Lines -= 1
		return
	}

	text.Runes[text.Lines] = text.Runes[text.Lines][:len(text.Runes[text.Lines])-1]
}

func (text *Content) ClearLine() {
	if text.Lines == 0 {
		return
	}

	text.Runes[text.Lines] = []rune{}
}

func (text *Content) NewLine() {
	text.Runes = append(text.Runes, []rune{})
	text.Lines += 1
}

func (text *Content) Show(screen tcell.Screen) {
	screen.Clear()
	for y, line := range text.Runes {
		for x, c := range line {
			screen.SetContent(x, y, c, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}
	}
	screen.Show()
}

func InitEit() tcell.Screen {
	tcell.SetEncodingFallback(tcell.EncodingFallbackUTF8)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault)

	s.Clear()

	return s
}

func main() {
	screen := InitEit()
	text := &Content{}

	events := make(chan tcell.Event)
	var event tcell.Event

	go func() {
		for {
			events <- screen.PollEvent()
		}
	}()

	for {
		text.Show(screen)
		select {
		case event = <-events:
		}

		switch e := event.(type) {
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyDEL:
				text.DeleteRune()
			case tcell.KeyEnter:
				text.NewLine()
			case tcell.KeyCtrlU:
				text.ClearLine()
			case tcell.KeyRune:
				text.AddRune(e.Rune())
			case tcell.KeyCtrlX, tcell.KeyCtrlC:
				screen.Fini()
				fmt.Fprintf(os.Stdout, "bye!\n")
				return
			}
		case *tcell.EventResize:
			screen.Show()
		case *tcell.EventInterrupt:
			os.Exit(0)
		}
	}
}
