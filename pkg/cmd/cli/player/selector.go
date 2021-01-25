package player

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Selector is panel for selecting music
type Selector struct {
	*tview.Flex
	app *App

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
			AddItem(s.djID, 0, 2, true).
			AddItem(s.musicTitle, 0, 3, false), 0, 1, false)

	s.SetKeyHandler()
	return s
}

func (s *Selector) SetKeyHandler() {
	s.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Rune() {
		case 'b':
			fmt.Println("fuag")
		}
		return e
	})
}
