package turntable

import (
	"sync"

	"github.com/rivo/tview"
)

// Panel is base for other functional panels
type Panel struct {
	once sync.Once
	tview.Flex
}

// Init initialize the panel
func (p *Panel) Init(app *App) {

}
