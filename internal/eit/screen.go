package eit

import (
	"fmt"
	"github.com/gdamore/tcell"
	"os"
)

func InitScreen() tcell.Screen {
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
