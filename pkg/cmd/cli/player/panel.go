package player

import (
	"sync"

	"github.com/rivo/tview"
)

// Panel is base for other functional panels
type Panel struct {
	once sync.Once
	tview.Flex

	app *App
}

// Init set app for the panel
func (p *Panel) Init(app *App) {
	p.once.Do(func() {
		p.app = app
	})
}

// App return pointer of the App
func (p *Panel) App() *App {
	return p.app
}

// NewPanel create and return the Panel instance
func NewPanel(app *App) *Panel {
	p := &Panel{}
	p.Init(app)
	return p
}
