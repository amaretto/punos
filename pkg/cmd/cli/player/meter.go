package player

import (
	"strings"

	"github.com/rivo/tview"
)

type MeterBox struct {
	*tview.Flex
	meters []*Meter
}

type Meter struct {
	*tview.TextView
	label string
}

func (mb *MeterBox) update() {
	for _, m := range mb.meters {
		val, max := 50, 100
		_, _, width, height := m.GetRect()
		m.SetText(m.genMeter(val, max, height, width, m.label))
	}
}

func NewMeterBox() *MeterBox {
	mb := &MeterBox{
		Flex:   tview.NewFlex(),
		meters: []*Meter{&Meter{tview.NewTextView(), "volume"}, &Meter{tview.NewTextView(), "speed"}, &Meter{tview.NewTextView(), "pitch"}, &Meter{tview.NewTextView(), "filter"}},
	}

	mb.SetDirection(tview.FlexColumn).SetBorder(true).SetTitle("Meters").SetTitleAlign(tview.AlignLeft).SetBorderPadding(2, 2, 2, 2)
	for _, m := range mb.meters {
		m.SetTextAlign(tview.AlignCenter)
		mb.AddItem(m, 0, 1, false)
	}

	mb.update()
	return mb
}

func (m *Meter) genMeter(val, max, height, width int, label string) string {
	var meterStr, guage, empty string
	// obtain space for header, footer, label
	border := (height - 3) * val / max
	if height > 6 && width > 9 {
		meterStr = "┌" + strings.Repeat("─", width-2) + "┐\n"
		guage = "│" + strings.Repeat("─", width-2) + "│\n"
		empty = "│" + strings.Repeat(" ", width-2) + "│\n"
		for i := height - 3; i > 0; i-- {
			if i > border {
				meterStr = meterStr + empty
			} else {
				meterStr = meterStr + guage
			}
		}
		meterStr = meterStr + "└" + strings.Repeat("─", width-2) + "┘\n" + label
	}
	return meterStr
}
