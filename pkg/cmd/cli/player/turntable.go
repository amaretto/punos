package player

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Turntable give some functions of music player
type Turntable struct {
	app *Player

	*tview.Flex
	helpModal    tview.Primitive
	isHelpActive bool

	djName        *DefaultView
	turntableID   *DefaultView
	musicTitle    *DefaultView
	progressBar   *ProgressBar
	waveformPanel *WaveformPanel
	playPauseBox  *PlayPausePanel
	meterBox      *MeterBox
}

func newTurntable(app *Player) *Turntable {
	t := &Turntable{
		app: app,

		Flex:      tview.NewFlex(),
		helpModal: newHelpModal(),

		djName:        NewDefaultView("DJ"),
		turntableID:   NewDefaultView("TurnTable"),
		musicTitle:    NewDefaultView("Music"),
		progressBar:   NewProgressBar(),
		waveformPanel: NewWaveformPanel(),
		playPauseBox:  NewPlayPausePanel(),
		meterBox:      NewMeterBox(),
	}

	t.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(t.turntableID, 0, 2, false).
			AddItem(t.djName, 0, 2, false).
			AddItem(t.musicTitle, 0, 3, false), 0, 1, false).
		AddItem(t.progressBar, 0, 1, false).
		AddItem(t.waveformPanel, 0, 6, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(t.playPauseBox, 0, 3, false).
			AddItem(t.meterBox, 0, 7, false), 0, 4, false)

	t.initTurntable()
	return t
}

func (t *Turntable) initTurntable() {
	// ToDo: set dj name and turntable from configuration or arguments
	t.djName.SetText(t.app.playerID)
	t.turntableID.SetText("TurnTable")
	t.SetKeyHandler()
}

func newHelpModal() tview.Primitive {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}

	rows, cols := 10, 2
	keyBindingHelp := [][]string{
		{"Key", "Description"},
		{"Space", "Play/Pause"},
		{"Esc", "Quit"},
		{"l", "Fast Forward"},
		{"h", "Rewind"},
		{"j", "Volume Up"},
		{"k", "Volume Down"},
		{"m", "Speed Up"},
		{",", "Speed Down"},
		{"c", "Set/Jump Cue"},
	}

	table := tview.NewTable().
		SetBorders(true)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
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

func (t *Turntable) update() {
	t.musicTitle.SetText(t.app.nowPlaying.Title)
	t.progressBar.update(t.app.streamer.Position(), t.app.streamer.Len())
	t.waveformPanel.update(t.app.nowPlaying.Waveform, t.app.streamer.Position())
	t.meterBox.update(int((t.app.volume.Volume+1)*100), int(t.app.resampler.Ratio()*100))
	// ToDo: update PlayPause
}

func (t *Turntable) SetKeyHandler() {
	t.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Rune() {
		case 'l':
			t.app.Fforward()
		case 'h':
			t.app.Rewind()
		case 'j':
			t.app.Voldown()
		case 'k':
			t.app.Volup()
		case 'm':
			t.app.Spdup()
		case ',':
			t.app.Spddown()
		case 'c':
			t.app.Cue()
		}
		return e
	})
}
