package player

import (
	"fmt"

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
	meterBox     *tview.Flex
}

// DefaultView is
type DefaultView struct {
	*tview.TextView
}

// NewDefaultView is
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

		djID:         NewDefaultView("DJ"),
		turntableID:  NewDefaultView("TurnTable"),
		musicTitle:   NewDefaultView("Music"),
		progressBar:  NewDefaultView("Progress"),
		waveformBox:  NewDefaultView("Waveform"),
		playPauseBox: NewDefaultView("Play/Pause"),
		meterBox:     tview.NewFlex(),
	}

	t.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(t.turntableID, 0, 2, false).
			AddItem(t.djID, 0, 2, false).
			AddItem(t.musicTitle, 0, 3, false), 0, 1, false).
		AddItem(t.progressBar, 0, 1, false).
		AddItem(t.waveformBox, 0, 6, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(t.playPauseBox, 0, 3, false).
			AddItem(t.meterBox, 0, 7, false), 0, 4, false)

	t.setKeyHandler()
	return t
}

func (t *Turntable) setKeyHandler() {
	t.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		t.app.SetGlobalKeyBinding(e)
		switch e.Key() {
		case tcell.KeyESC:
			t.app.Stop()
		}

		switch e.Rune() {
		case 'a':
			fmt.Println(t.HasFocus())
		}
		return e
	})
}

func (t *Turntable) update() {
	// get music info from app and display it
}
