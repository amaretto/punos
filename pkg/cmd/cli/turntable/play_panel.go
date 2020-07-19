package turntable

import (
	"unicode"

	"github.com/gdamore/tcell"
)

// PlayPanel give some functions of music player
type PlayPanel struct {
	Panel
}

// HandleEvent handles key event
func (p *PlayPanel) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEsc:
			p.App().Stop()
			return true
		}

		switch unicode.ToLower(ev.Rune()) {
		case ' ':
			p.App().Stop()
			return true
		}

	}
	return true
}
