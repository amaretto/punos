package player

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell/v2"
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

	// Music
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	sampleRate beep.SampleRate
	resampler  *beep.Resampler
	volume     *effects.Volume
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

	a.s = newSelector(a)
	a.pages.AddPage("selector", a.s, true, false)

	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			a.Stop()
		}

		switch event.Rune() {
		case 'n':
			a.pages.SwitchToPage("selector")
			a.app.SetFocus(a.s)
		case 'f':
			a.pages.SwitchToPage("turntable")
			a.app.SetFocus(a.t)
		}
		return event
	})

	return a
}

// Start kick the application
func (a *App) Start() {
	a.LoadMusic("test")
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

func (a *App) SetGlobalKeyBinding(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyESC:
		a.Stop()
	}

	switch event.Rune() {
	case 'n':
		a.pages.SwitchToPage("selector")
		a.app.SetFocus(a.s)
	case 'f':
		a.pages.SwitchToPage("turntable")
		a.app.SetFocus(a.t)
	}
}

func (a *App) LoadMusic(path string) {
	//speaker.Lock()
	//a.ctrl.Paused = true
	//speaker.Unlock()
	//a.streamer.Close()

	f, err := os.Open("mp3/test.mp3")
	if err != nil {
		report(err)
	}
	// update title
	a.musicTitle = path

	//var format beep.Format
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		report(err)
	}
	// ToDo close streamer when music is switched

	//	a.wave.Wave = LoadWave(a.wave.WaveDirPath, a.musicTitle)
	//	a.wave.NormalizeWave()
	a.sampleRate = format.SampleRate
	a.streamer = streamer
	a.ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, a.streamer)}
	a.resampler = beep.ResampleRatio(4, 1, a.ctrl)
	a.volume = &effects.Volume{Streamer: a.resampler, Base: 2}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))
	speaker.Play(a.volume)

	//speaker.Lock()
	//a.ctrl.Paused = !a.ctrl.Paused
	//speaker.Unlock()
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
