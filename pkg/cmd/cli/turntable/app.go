package turntable

import "github.com/rivo/tview"

// App is standalone turntable(music player) application
type App struct {
	app *tview.Application

	pp *PlayPanel
	sp *SelectPanel
}

// Stop stop the application
func (a *App) Stop() {
	a.app.Stop()
}
