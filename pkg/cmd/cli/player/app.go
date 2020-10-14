package player

import (
	"strconv"
	"time"

	"github.com/rivo/tview"
)

// App is standalone dj player application
type App struct {
	app      *tview.Application
	playerID string
	t        *Turntable
	s        *Selector

	pages      *tview.Pages
	musicTitle string
	musicPath  string
	isPlay     bool
}

// New return App instance
func New() *App {
	a := &App{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),
	}

	a.t = newTurntable(a)
	a.pages.AddPage("turntable", a.t, true, true)
	a.pages.SwitchToPage("turntable")

	return a
}

// Start kick the application
func (a *App) Start() {
	go func() {
		for {
			time.Sleep(1 * time.Millisecond)
			a.t.musicTitle.SetText(strconv.FormatInt(time.Now().UnixNano(), 10))
			a.app.Draw()
		}
	}()
	if err := a.app.SetRoot(a.pages, true).Run(); err != nil {
		panic(err)
	}
}

// Stop stop the application
func (a *App) Stop() {
	a.app.Stop()
}

//	waveformString := `
//________   ___  ___   ________    ________   ________
//|\   __ \ |\  \|\  \ |\   ___  \ |\   __  \ |\   ____\
//\ \  \|\ \ \ \  \\  \\ \  \\ \  \\ \  \|\  \\ \  \___|_
// \ \   ___\ \ \  \\  \\ \  \\ \  \\ \  \\\  \\ \_____  \
//  \ \  \___| \ \  \\  \\ \  \\ \  \\ \  \\\  \\|____|\  \
//   \ \__\     \ \_______\\\__\\ \__\\ \_______\ ____\_\  \
//    \|__|      \|_______| \|__| \|__| \|_______||\________\
//                                                \|________|
//`
//	meterBox := tview.NewFlex()
//	meterBox.SetDirection(tview.FlexColumn).SetBorder(true).SetTitle("Meters").SetTitleAlign(tview.AlignLeft)
//	volumeMeter := tview.NewTextView().SetText("┌─────┐\n│     │\n│     │\n│     │\n│     │\n│─025─│\n│─────│\n│─────│\n│─────│\n└─────┘\nspeed").SetTextAlign(tview.AlignCenter)
//	volumeMeter2 := tview.NewTextView().SetText("┌─────┐\n│     │\n│     │\n│     │\n│     │\n│-010-│\n│     │\n│     │\n│─────│\n└─────┘\nvolume").SetTextAlign(tview.AlignCenter)
//	volumeMeter3 := tview.NewTextView().SetText("┌─────┐\n│     │\n│     │\n│     │\n│─────│\n│─120-│\n│─────│\n│─────│\n│─────│\n└─────┘\nbpm").SetTextAlign(tview.AlignCenter)
//	volumeMeter4 := tview.NewTextView().SetText("┌─────┐\n│     │\n│─────│\n│─────│\n│─────│\n│─080─│\n│─────│\n│─────│\n│─────│\n└─────┘\nfilter").SetTextAlign(tview.AlignCenter)
//	meterBox.AddItem(volumeMeter, 0, 1, false).AddItem(volumeMeter2, 0, 1, false).AddItem(volumeMeter3, 0, 1, false).AddItem(volumeMeter4, 0, 1, false)
//
//	//meterBox := tview.NewBox().SetBorder(true).SetTitle("Meters").SetTitleAlign(tview.AlignLeft)
//
//
//	dummyPage := tview.NewTextView()
//	dummyPage.SetText("hogehogehoge").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow)
//
//	pages := tview.NewPages()
//	pages.AddPage("ttpanel", flex, true, true)
//	pages.AddPage("dummyPage", dummyPage, true, false)
//
//	pages.SwitchToPage("ttpanel")
//	go func() {
//		for {
//			time.Sleep(1 * time.Millisecond)
//			musicTitle.SetText(strconv.FormatInt(time.Now().UnixNano(), 10))
//			app.Draw()
//			//pages.SwitchToPage("ttpanel")
//			//time.Sleep(1 * time.Second)
//			//pages.SwitchToPage("dummyPage")
//		}
//	}()
//
//	if err := app.SetRoot(pages, true).Run(); err != nil {
//		panic(err)
//	}
//
//}
