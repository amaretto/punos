package ui

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// TitleBar s
type TitleBar struct {
	once sync.Once
	views.SimpleStyledTextBar
}

// Init is
func (tb *TitleBar) Init() {
	tb.once.Do(func() {
		normal := tcell.StyleDefault.
			Foreground(tcell.ColorBlack).
			Background(tcell.ColorSilver)
		alternate := tcell.StyleDefault.
			Foreground(tcell.ColorBlue).
			Background(tcell.ColorSilver)

		tb.SimpleStyledTextBar.Init()
		tb.SimpleStyledTextBar.SetStyle(normal)
		tb.RegisterLeftStyle('N', normal)
		tb.RegisterLeftStyle('A', alternate)
		tb.RegisterCenterStyle('N', normal)
		tb.RegisterCenterStyle('A', alternate)
		tb.RegisterRightStyle('N', normal)
		tb.RegisterRightStyle('A', alternate)
	})

}

// NewTitleBar is
func NewTitleBar() *TitleBar {
	tb := &TitleBar{}
	tb.Init()
	return tb
}
