package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/x1ah/eit/internal/eit"
	"os"
)

func main() {
	screen := eit.InitScreen()
	config := &eit.Config{}
	buffer := eit.NewBuffer(config)

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
