package player

import "github.com/rivo/tview"

// Selector is panel for selecting music
type Selector struct {
	app *App
	*tview.Flex
}

func newSelector(app *App) *Selector {
	s := &Selector{
		app:  app,
		Flex: tview.NewFlex(),
	}
	return s
}
