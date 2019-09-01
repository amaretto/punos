package ui

import (
	"sync"

	"github.com/gdamore/tcell/views"
)

// Panel is
type Panel struct {
	tb *TitleBar
	//	sb   *StatusBar
	kb   *KeyBar
	once sync.Once
	app  *App

	views.Panel
}

// SetTitle is
func (p *Panel) SetTitle(title string) {
	p.tb.SetCenter(title)
}

// SetKeys is
func (p *Panel) SetKeys(words []string) {
	p.kb.SetKeys(words)
}

//// SetTitle is
//func (p *Panel) SetTitle(title string) {
//	p.tb.SetCenter(title)
//}

// Init is
func (p *Panel) Init(app *App) {
	p.once.Do(func() {
		p.app = app

		// create and set title bar
		p.tb = NewTitleBar()
		p.tb.SetRight(app.GetAppName())
		p.tb.SetCenter(" ")

		p.Panel.SetTitle(p.tb)

		// create and set key bar
		p.kb = NewKeyBar()
		p.kb.SetCenter(" ")

		p.Panel.SetStatus(p.kb)
	})
}

// App is
func (p *Panel) App() *App {
	return p.app
}

// NewPanel is
func NewPanel(app *App) *Panel {
	p := &Panel{}
	p.Init(app)
	return p
}
