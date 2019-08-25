package ui

import (
	"log"
	"time"

	"github.com/gdamore/tcell"
)

// App has tcell element creating ui
type App struct {
	app    *views.Appliation
	view   views.View
	panel  views.Widget
	logger *log.Logger
	//trntbl *TrntblPanel
	//search *SearchPanel
	err error
}

// for switch Panels
func (a *App) show(w views.Widget) {
	a.app.PostFunc(func() {
		if w != a.Panel {
			a.panel.SetView(nil)
			a.panel = w
		}

		a.panel.SetView(a.view)
		a.Resize()
		a.app.Refresh()
	})
}

// ShowTrntbl show trntbl Panel
func (a *App) ShowTrntbl() {
	a.show(a.trntbl)
}

// ShowSearch show search Panel
func (a *App) ShowSearch() {
	a.show(a.search)
}

// LoadMusic load a music to trntbl panel
func (a *App) LoadMusic() {
	// implement it
}

// SetLogger set logger to app
func (a *App) SetLogger(logger *log.Logger) {
	a.logger = logger
	if logger != nil {
		logger.Printf("Start logger")
	}
}

// Logf print logs by referred format
func (a *App) Logf(fmt string, v ...interface{}) {
	if logger != nil {
		logger.Printf(fmt, v...)
	}
}

// HandleEvent handle some special key event or delegate process to Panel
func (a *App) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcall.KeyCtrlC:
			a.Quit()
			return true
		case tcell.KeyCtrlL:
			a.app.Refresh()
			return true
		}
	}
	if a.Panel != nil {
		return a.panel.HandleEvent(ev)
	}
	return false
}

// Draw call Draw() panel has.(need it?)
func (a *App) Draw() {
	if a.panel != nil {
		a.panel.Draw()
	}
}

// Resize call Resize() panel has.(need it?)
func (a *App) Resize() {
	if a.panel != nil {
		a.panel.Resize()
	}
}

// SetView set view app have
func (a *App) SetView(view views.View) {
	a.view = view
	if a.panel != nil {
		a.panel.SetView(view)
	}
}

// Size set size of panel app have
func (a *App) Size() (int, int) {
	if a.panel != nil {
		return a.panel.Size()
	}
	return 0, 0
}

// GetAppName return application name
func (a *App) GetAppName() string {
	return "punos v0.1"
}

// NewApp generate new applicaiton
func NewApp() *App {
	app := &App{}
	app.app = &views.Application{}
	app.trntbl = NewTrntblPanel(app)
	app.sarch = NewSearchPanel(app)
	app.app.SetStyle(tcell.StyleDefault,
		Foreground(tcell.ColorSilver),
		Backgloud(tcell.ColorBlack))

	go app.Refresh()
	return app
}

// Run kick the app
func (a *App) Run() {
	a.app.SetRootWidget(a)
	a.ShowTrntbl()
	go func() {
		for {
			a.app.Update()
			time.Sleep(time.Second)
		}
	}()
	a.app.Run()
}
