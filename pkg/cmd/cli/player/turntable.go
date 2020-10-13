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
	djID         *DefaultView
	turntableID  *DefaultView
	musicTitle   *DefaultView
	progressBar  *DefaultView
	waveformBox  *DefaultView
	playPauseBox *DefaultView
	meterBox     tview.Flex
}

type DefaultView struct {
	*tview.TextView
}

func NewDefaultView(title string) *DefaultView {
	d := &DefaultView{
		TextView: tview.NewTextView(),
	}
	d.TextView.SetBorder(true).SetTitleAlign(tview.AlignLeft).SetTitle(title)
	d.TextView.SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow)
	return d
}

func newTurntable(app *App) *Turntable {
	t := &Turntable{
		app:  app,
		Flex: tview.NewFlex(),

		turntableID:  NewDefaultView("TurnTable"),
		musicTitle:   NewDefaultView("Music"),
		progressBar:  NewDefaultView("Progress"),
		waveformBox:  NewDefaultView("Waveform"),
		playPauseBox: NewDefaultView("Play/Pause"),
		meterBox:     tview.NewFlex(),
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
