package server

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/spf13/cobra"
)

// NewCommand create command
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "server",
		Short: "start punos server",
		Long:  `start punos server`,
		Run: func(cmd *cobra.Command, args []string) {
			callBox()
		},
	}
	return c
}

func callBox() {
	app := tview.NewApplication()

	// turntable id
	turntableID := tview.NewTextView()
	turntableID.SetBorder(true).SetTitle("Server").SetTitleAlign(tview.AlignLeft)
	turntableID.SetText("Server 1").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow)

	// dj id
	djID := tview.NewTextView()
	djID.SetBorder(true).SetTitle("DJ").SetTitleAlign(tview.AlignLeft)
	djID.SetText("amaretto").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow)

	// music title
	musicTitle := tview.NewTextView()
	musicTitle.SetBorder(true).SetTitle("Music").SetTitleAlign(tview.AlignLeft)
	musicTitle.SetText("WHAT YOU GOT").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow)

	progressBar := tview.NewTextView()
	progressBar.SetBorder(true).SetTitle("Progress").SetTitleAlign(tview.AlignLeft)
	progressBar.SetText("[=========================>------------------] 3m15s/4m11s").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow)

	waveformBox := tview.NewTextView()
	waveformBox.SetBorder(true).SetTitle("Waveform").SetTitleAlign(tview.AlignLeft)

	waveformString := `
________   ___  ___   ________    ________   ________      
|\   __ \ |\  \|\  \ |\   ___  \ |\   __  \ |\   ____\     
\ \  \|\ \ \ \  \\  \\ \  \\ \  \\ \  \|\  \\ \  \___|_    
 \ \   ___\ \ \  \\  \\ \  \\ \  \\ \  \\\  \\ \_____  \   
  \ \  \___| \ \  \\  \\ \  \\ \  \\ \  \\\  \\|____|\  \  
   \ \__\     \ \_______\\\__\\ \__\\ \_______\ ____\_\  \ 
    \|__|      \|_______| \|__| \|__| \|_______||\________\
                                                \|________|
`
	waveformBox.SetText(waveformString).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow)

	playpauseBox := tview.NewBox().SetBorder(true).SetTitle("PlayPause").SetTitleAlign(tview.AlignLeft)

	meterBox := tview.NewFlex()
	meterBox.SetDirection(tview.FlexColumn).SetBorder(true).SetTitle("Meters").SetTitleAlign(tview.AlignLeft)
	volumeMeter := tview.NewTextView().SetText("┌─────┐\n│     │\n│     │\n│     │\n│     │\n│─025─│\n│─────│\n│─────│\n│─────│\n└─────┘\nspeed").SetTextAlign(tview.AlignCenter)
	volumeMeter2 := tview.NewTextView().SetText("┌─────┐\n│     │\n│     │\n│     │\n│     │\n│-010-│\n│     │\n│     │\n│─────│\n└─────┘\nvolume").SetTextAlign(tview.AlignCenter)
	volumeMeter3 := tview.NewTextView().SetText("┌─────┐\n│     │\n│     │\n│     │\n│─────│\n│─120-│\n│─────│\n│─────│\n│─────│\n└─────┘\nbpm").SetTextAlign(tview.AlignCenter)
	volumeMeter4 := tview.NewTextView().SetText("┌─────┐\n│     │\n│─────│\n│─────│\n│─────│\n│─080─│\n│─────│\n│─────│\n│─────│\n└─────┘\nfilter").SetTextAlign(tview.AlignCenter)
	meterBox.AddItem(volumeMeter, 0, 1, false).AddItem(volumeMeter2, 0, 1, false).AddItem(volumeMeter3, 0, 1, false).AddItem(volumeMeter4, 0, 1, false)

	//meterBox := tview.NewBox().SetBorder(true).SetTitle("Meters").SetTitleAlign(tview.AlignLeft)

	// layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(turntableID, 0, 2, false).
			AddItem(djID, 0, 2, false).
			AddItem(musicTitle, 0, 3, false), 0, 1, false).
		AddItem(progressBar, 0, 1, false).
		AddItem(waveformBox, 0, 6, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(playpauseBox, 0, 3, false).
			AddItem(meterBox, 0, 7, false), 0, 4, false)

	dummyPage := tview.NewTextView()
	dummyPage.SetText("hogehogehoge").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow)

	pages := tview.NewPages()
	pages.AddPage("ttpanel", flex, true, true)
	pages.AddPage("dummyPage", dummyPage, true, false)

	pages.SwitchToPage("ttpanel")
	go func() {
		for {
			time.Sleep(1 * time.Millisecond)
			musicTitle.SetText(strconv.FormatInt(time.Now().UnixNano(), 10))
			app.Draw()
			//pages.SwitchToPage("ttpanel")
			//time.Sleep(1 * time.Second)
			//pages.SwitchToPage("dummyPage")
		}
	}()

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}

}
