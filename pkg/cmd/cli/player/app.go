package player

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Player is standalone dj player application
type Player struct {
	app       *tview.Application
	playerID  string
	pages     *tview.Pages
	turntable *Turntable
	selector  *Selector

	// music info
	musicTitle string
	musicPath  string

	// music
	isPlay     bool
	streamer   beep.StreamSeekCloser
	ctrl       *beep.Ctrl
	sampleRate beep.SampleRate
	resampler  *beep.Resampler
	volume     *effects.Volume

	// waveform
	wf []int
}

// New return App instance
func New() *Player {
	p := &Player{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),
	}

	p.turntable = newTurntable(p)
	p.pages.AddPage("turntable", p.turntable, true, true)
	p.pages.SwitchToPage("turntable")

	p.selector = newSelector(p)
	p.pages.AddPage("selector", p.selector, true, false)

	p.setAppGlobalKeyBinding()
	// need to set focus to Primitive in Flex(HasFocus)
	p.app.SetFocus(p.turntable.waveformPanel)

	// init speaker
	// source : speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))
	var sampleRate beep.SampleRate = 44100
	speaker.Init(sampleRate, int(time.Duration(sampleRate)*(time.Second/30)/time.Second))

	return p
}

func (p *Player) setAppGlobalKeyBinding() {
	p.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		// global key binding
		switch event.Key() {
		case tcell.KeyESC:
			p.Stop()
		}
		switch event.Rune() {
		case ' ':
			speaker.Lock()
			p.ctrl.Paused = !p.ctrl.Paused
			speaker.Unlock()
		case 's':
			p.pages.SwitchToPage("selector")
			p.app.SetFocus(p.selector.musicListView)
		case 't':
			p.pages.SwitchToPage("turntable")
			p.app.SetFocus(p.turntable.waveformPanel)
		}

		return event
	})
}

func (p *Player) Start() {
	p.LoadMusic("mp3/test.mp3")
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
			p.turntable.musicTitle.SetText(p.musicTitle)
			p.app.Draw()
			p.turntable.progressBar.update(p.streamer.Position(), p.streamer.Len())
		}
	}()
	if err := p.app.SetRoot(p.pages, true).Run(); err != nil {
		panic(err)
	}
}

// Stop stop the application
func (p *Player) Stop() {
	p.app.Stop()
}

func (p *Player) LoadMusic(path string) {
	if p.ctrl != nil {
		speaker.Lock()
		p.ctrl.Paused = true
		speaker.Unlock()
		p.streamer.Close()
	}

	f, err := os.Open(path)
	if err != nil {
		report(err)
	}
	// update title
	p.musicTitle = path

	//var format beep.Format
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		report(err)
	}
	// ToDo close streamer when music is switched

	p.sampleRate = format.SampleRate
	p.streamer = streamer
	p.ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, p.streamer)}
	p.resampler = beep.ResampleRatio(4, 1, p.ctrl)
	p.volume = &effects.Volume{Streamer: p.resampler, Base: 2}

	speaker.Play(p.volume)

	speaker.Lock()
	p.ctrl.Paused = !p.ctrl.Paused
	speaker.Unlock()
}

func loadWaveform(title string) {
	dbPath := "mp3/test.db"

	// ToDo: Implement error handling
	con, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalln(err)
	}

	cmd := "SELECT wave FROM waveform WHERE title = ?"

	_, err = con.Exec(cmd, title)
	if err != nil {
		log.Fatalln(err)
	}
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
