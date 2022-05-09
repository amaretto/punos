package player

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newHelpModal(keyBindingHelp [][]string) tview.Primitive {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}

	table := tview.NewTable().
		SetBorders(true)

	for r := 0; r < len(keyBindingHelp); r++ {
		for c := 0; c < 2; c++ {
			color := tcell.ColorWhite
			if c < 1 || r < 1 {
				color = tcell.ColorYellow
			}
			table.SetCell(r, c,
				tview.NewTableCell(keyBindingHelp[r][c]).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
		}
	}

	return modal(table, calcWidth(keyBindingHelp), len(keyBindingHelp)*2+1)
}

// for resize width
func calcWidth(table [][]string) int {
	var maxKeyLen, maxDescLen int
	for i := 0; i < len(table); i++ {
		if len(table[i][0]) > maxKeyLen {
			maxKeyLen = len(table[i][0])
		}
		if len(table[i][1]) > maxDescLen {
			maxDescLen = len(table[i][1])
		}
	}
	// 3 = boarder count
	return maxKeyLen + maxDescLen + 3
}
