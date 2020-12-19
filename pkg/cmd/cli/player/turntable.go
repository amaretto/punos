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
	djName       *DefaultView
	turntableID  *DefaultView
	musicTitle   *DefaultView
	progressBar  *DefaultView
	waveformPanel  *WaveformPanel
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

		djName:       NewDefaultView("DJ"),
		turntableID:  NewDefaultView("TurnTable"),
		musicTitle:   NewDefaultView("Music"),
		progressBar:  NewDefaultView("Progress"),
		waveformPanel:  NewWaveformPanel(),
		playPauseBox: NewDefaultView("Play/Pause"),
		meterBox:     tview.NewFlex(),
	}

	t.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(t.turntableID, 0, 2, false).
			AddItem(t.djName, 0, 2, false).
			AddItem(t.musicTitle, 0, 3, false), 0, 1, false).
		AddItem(t.progressBar, 0, 1, false).
		AddItem(t.waveformPanel, 0, 6, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(t.playPauseBox, 0, 3, false).
			AddItem(t.meterBox, 0, 7, false), 0, 4, false)

	t.initTurntable()
	t.SetKeyHandler()
	return t
}

func (t *Turntable) initTurntable() {
	// ToDo: set dj name and turntable from configuration or arguments
	t.djName.SetText("anonymous")
	t.turntableID.SetText("TurnTable")
}

func (t *Turntable) SetKeyHandler() {
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
