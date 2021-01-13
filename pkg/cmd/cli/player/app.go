package player

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// App is standalone dj player application
type App struct {
	app      *tview.Application
	playerID string
	t        *Turntable
	s        *Selector

	pages      *tview.Pages
	musicTitle string
	musicPath  string
	isPlay     bool
}

// New return App instance
func New() *App {
	a := &App{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),
	}

	a.t = newTurntable(a)
	a.pages.AddPage("turntable", a.t, true, true)
	a.pages.SwitchToPage("turntable")

	a.s = newSelector(a)
	a.pages.AddPage("selector", a.s, true, false)

	return a
}

// Start kick the application
func (a *App) Start() {
	go func() {
		for {
			time.Sleep(1 * time.Millisecond)
			a.t.musicTitle.SetText(strconv.FormatInt(time.Now().UnixNano(), 10))
			a.app.Draw()
		}
	}()
	if err := a.app.SetRoot(a.pages, true).Run(); err != nil {
		panic(err)
	}
}

// Stop stop the application
func (a *App) Stop() {
	a.app.Stop()
}

func (a *App) SetGlobalKeyBinding(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyESC:
		a.Stop()
	}

	switch event.Rune() {
	case 'n':
		a.pages.SwitchToPage("selector")
		a.app.SetFocus(a.s)
	case 'f':
		a.pages.SwitchToPage("turntable")
		a.app.SetFocus(a.t)
	}
}
