package player

import "github.com/rivo/tview"

type MeterBox struct {
	*tview.Flex
}

type Meter struct {
	Max   float32
	Min   float32
	Scale float32
}

func NewMeterBox() *MeterBox {
	m := &MeterBox{
		Flex: tview.NewFlex(),
	}

	m.SetDirection(tview.FlexColumn).SetBorder(true).SetTitle("Meters").SetTitleAlign(tview.AlignLeft)

	volumeMeter := tview.NewTextView().SetText("┌─────┐\n│     │\n│     │\n│     │\n│     │\n│─025─│\n│─────│\n│─────│\n│─────│\n└─────┘\nspeed").SetTextAlign(tview.AlignCenter)
	volumeMeter2 := tview.NewTextView().SetText("┌─────┐\n│     │\n│     │\n│     │\n│     │\n│-010-│\n│     │\n│     │\n│─────│\n└─────┘\nvolume").SetTextAlign(tview.AlignCenter)
	volumeMeter3 := tview.NewTextView().SetText("┌─────┐\n│     │\n│     │\n│     │\n│─────│\n│─120-│\n│─────│\n│─────│\n│─────│\n└─────┘\nbpm").SetTextAlign(tview.AlignCenter)
	volumeMeter4 := tview.NewTextView().SetText("┌─────┐\n│     │\n│─────│\n│─────│\n│─────│\n│─080─│\n│─────│\n│─────│\n│─────│\n└─────┘\nfilter").SetTextAlign(tview.AlignCenter)
	m.AddItem(volumeMeter, 0, 1, false).AddItem(volumeMeter2, 0, 1, false).AddItem(volumeMeter3, 0, 1, false).AddItem(volumeMeter4, 0, 1, false)
	return m
}
