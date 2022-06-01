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

func (mb *MeterBox) update(volume, speed int) {
	var val, max int
	for i, m := range mb.meters {

		if i == 0 {
			val, max = volume, 200
		} else if i == 1 {
			val, max = speed, 200
		} else {
			val, max = 100, 200
		}
		if val > 200 {
			val = 200
		}
		_, _, width, height := m.GetRect()
		if width >= 15 {
			width = 15
		}
		m.SetText(m.genMeter(val, max, height, width, m.label))
	}
}

func NewMeterBox() *MeterBox {
	mb := &MeterBox{
		Flex:   tview.NewFlex(),
		meters: []*Meter{&Meter{tview.NewTextView(), "volume"}, &Meter{tview.NewTextView(), "speed"}, &Meter{tview.NewTextView(), "N/A"}, &Meter{tview.NewTextView(), "N/A"}},
	}

	mb.SetDirection(tview.FlexColumn).SetBorder(true).SetTitle("Meters").SetTitleAlign(tview.AlignLeft).SetBorderPadding(2, 2, 2, 2)
	for _, m := range mb.meters {
		m.SetTextAlign(tview.AlignCenter)
		mb.AddItem(m, 0, 1, false)
	}

	mb.update(100, 200)
	return mb
}

func (m *Meter) genMeter(val, max, height, width int, label string) string {
	var meterStr, guage, empty string
	// obtain space for header, footer, label
	border := (height - 3) * 8 * val / max
	remain := border % 8
	isRemain := true
	if height > 6 && width > 9 {
		meterStr = "┌" + strings.Repeat("─", width-2) + "┐\n"
		guage = "│" + strings.Repeat("█", width-2) + "│\n"
		empty = "│" + strings.Repeat(" ", width-2) + "│\n"
		for i := (height - 3) * 8; i > 0; i -= 8 {
			if i > border {
				meterStr = meterStr + empty
			} else if i <= border && isRemain {
				var c string
				switch remain {
				case 0:
					c = " "
				case 1:
					c = "▁"
				case 2:
					c = "▂"
				case 3:
					c = "▃"
				case 4:
					c = "▄"
				case 5:
					c = "▅"
				case 6:
					c = "▆"
				case 7:
					c = "▇"
				}
				tmpGauge := "│" + strings.Repeat(c, width-2) + "│\n"
				meterStr = meterStr + tmpGauge
				isRemain = false
			} else {
				meterStr = meterStr + guage
			}
		}
		meterStr = meterStr + "└" + strings.Repeat("─", width-2) + "┘\n" + label
	}
	return meterStr
}
