package player

import "github.com/rivo/tview"

// App is standalone dj player application
type App struct {
	app      *tview.Application
	playerID string
	t        *Turntable
	s        *Selector

	musicTitle string
	musicPath  string
	isPlay     bool
}

func (a *App) Start() {

}

// Stop stop the application
func (a *App) Stop() {
	a.app.Stop()
}
