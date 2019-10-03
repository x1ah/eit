package cmd

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/x1ah/eit/internal/eit"
	"os"
)

func Run() {
	screen := eit.InitScreen()
	filepath := ""
	if len(os.Args) >= 2 {
		filepath = os.Args[1]
	}
	config := &eit.Config{FilePath: filepath}
	buffer, err := eit.NewBuffer(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}

	events := make(chan tcell.Event)
	var event tcell.Event

	go func() {
		for {
			events <- screen.PollEvent()
		}
	}()

	for {
		buffer.Draw(screen)
		select {
		case event = <-events:
		}

		switch e := event.(type) {
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyDEL:
				buffer.Delete()
			case tcell.KeyEnter:
				buffer.NewLine()
			case tcell.KeyRune:
				buffer.InsertRune(e.Rune())
			case tcell.KeyLeft:
				buffer.CurrCursor.MoveLeft()
			case tcell.KeyRight:
				buffer.CurrCursor.MoveRight()
			case tcell.KeyUp:
				buffer.CurrCursor.MoveUp()
			case tcell.KeyDown:
				buffer.CurrCursor.MoveDown()
			case tcell.KeyCtrlS:
				buffer.SaveAs()
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
