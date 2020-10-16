package player

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Selector is panel for selecting music
type Selector struct {
	app *App
	*tview.Flex
}

func newSelector(app *App) *Selector {
	s := &Selector{
		app:  app,
		Flex: tview.NewFlex(),
	}
	s.SetTitle("selector")

	s.setKeyHandler()
	return s
}

func (s *Selector) setKeyHandler() {
	s.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Key() {
		case tcell.KeyESC:
			s.app.Stop()
		}
		switch e.Rune() {
		case 'n':
			s.app.pages.SwitchToPage("turntable")
		}
		return e
	})
}
