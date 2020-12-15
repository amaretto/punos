package player

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Selector is panel for selecting music
type Selector struct {
	app *App
	*tview.Flex

	djID        *DefaultView
	turntableID *DefaultView
	musicTitle  *DefaultView
}

func newSelector(app *App) *Selector {
	s := &Selector{
		app:  app,
		Flex: tview.NewFlex(),

		djID:        NewDefaultView("DJ"),
		turntableID: NewDefaultView("TurnTable"),
		musicTitle:  NewDefaultView("Music"),
	}
	s.SetTitle("selector")

	s.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(s.turntableID, 0, 2, false).
			AddItem(s.djID, 0, 2, false).
			AddItem(s.musicTitle, 0, 3, false), 0, 1, false)

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
			fmt.Println("hoge")
			s.app.pages.SwitchToPage("turntable")
		}
		return e
	})
}
