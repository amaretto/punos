package player

import (
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Turntable give some functions of music player
type Turntable struct {
	app *App
	*tview.Flex

	// CHANGE IT
	djID         tview.TextView
	musicTitle   tview.TextView
	progressBar  tview.TextView
	waveformBox  tview.TextView
	playPauseBox tview.TextView
	meterBox     tview.Flex
}

func newTurntable(app *App) *Turntable {
	t := &Turntable{
		app:  app,
		Flex: tview.NewFlex(),
	}
	return t
}

// HandleEvent handles key event of Turntable
func (t *Turntable) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEsc:
			t.app.Stop()
			return true
		}

		switch unicode.ToLower(ev.Rune()) {
		case ' ':
			t.app.Stop()
			return true
		}
	}
	return true
}

func (t *Turntable) update() {
	// get music info from app and display it
}
