package player

import (
	"unicode"

	"github.com/gdamore/tcell"
)

// Turntable give some functions of music player
type Turntable struct {
	Panel
}

// HandleEvent handles key event of Turntable
func (p *Turntable) HandleEvent(ev tcell.Event) bool {
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

func (t *Turntable) update() {
	// get music info from app and display it
}
