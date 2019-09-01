package ui

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// KeyBar is
type KeyBar struct {
	once sync.Once
	views.SimpleStyledTextBar
}

// Init is
func (k *KeyBar) Init() {
	k.once.Do(func() {
		normal := tcell.StyleDefault.
			Foreground(tcell.ColorBlack).
			Background(tcell.ColorSilver)
		alternate := tcell.StyleDefault.
			Foreground(tcell.ColorBlue).
			Background(tcell.ColorSilver).Bold(true)

		k.SimpleStyledTextBar.Init()
		k.SimpleStyledTextBar.SetStyle(normal)
		k.RegisterLeftStyle('N', normal)
		k.RegisterLeftStyle('A', alternate)
	})
}

// SetKeys is
func (k *KeyBar) SetKeys(words []string) {
	b := make([]rune, 0, 80)
	for i, w := range words {
		esc := false
		if i != 0 && len(w) != 0 {
			b = append(b, ' ')
		}
		for _, r := range w {
			if esc {
				if r == ']' {
					b = append(b, '%', 'N')
					esc = false
				} else if r == '%' {
					b = append(b, '%')
				}
				b = append(b, r)

			} else {
				b = append(b, r)
				if r == '[' {
					esc = true
					b = append(b, '%', 'A')
				} else if r == '%' {
					b = append(b, '%')
				}
			}
		}
	}
	k.SetLeft(string(b))
}

// NewKeyBar is
func NewKeyBar() *KeyBar {
	kb := &KeyBar{}
	kb.Init()
	return kb
}
