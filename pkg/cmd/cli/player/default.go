package player

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// DefaultView is simple box view which only show some text
type DefaultView struct {
	*tview.TextView
}

func NewDefaultView(title string) *DefaultView {
	d := &DefaultView{
		TextView: tview.NewTextView(),
	}
	d.TextView.SetBorder(true).SetTitleAlign(tview.AlignLeft).SetTitle(title)
	d.TextView.SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow)
	return d
}
