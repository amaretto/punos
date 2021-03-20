package player

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
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
	wf []byte

	// cue
	cuePoint int
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
			p.turntable.waveformPanel.update(p.wf, p.streamer.Position())
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

	p.loadWaveform(path)
}

func (p *Player) loadWaveform(path string) {
	dbPath := "mp3/test.db"

	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}

	cmd := "SELECT wave FROM waveform WHERE path = ?"

	row := db.QueryRow(cmd, filepath.Base(path))
	if err != nil {
		log.Fatalln(err)
	}
	var data []byte
	row.Scan(&data)
	logrus.Debug(data)
	p.wf = data
}

// Fforward fast-forward music
func (p *Player) Fforward() {
	speaker.Lock()
	newPos := p.streamer.Position() + p.sampleRate.N(time.Millisecond*100)
	if newPos >= p.streamer.Len() {
		newPos = p.streamer.Len() - 1
	}
	if err := p.streamer.Seek(newPos); err != nil {
		report(err)
	}
	speaker.Unlock()
}

// Rewind rewind music
func (p *Player) Rewind() {
	speaker.Lock()
	newPos := p.streamer.Position() - p.sampleRate.N(time.Millisecond*100)
	if newPos < 0 {
		newPos = 0
	}
	if err := p.streamer.Seek(newPos); err != nil {
		report(err)
	}
	speaker.Unlock()
}

// Cue set and return cue point
func (p *Player) Cue() {
	// ToDo : adopt multiple cue point
	if p.ctrl.Paused {
		speaker.Lock()
		p.cuePoint = p.streamer.Position()
		speaker.Unlock()
	} else {
		speaker.Lock()
		p.streamer.Seek(p.cuePoint)
		speaker.Unlock()
	}
}

// Volup increase volume of music
func (p *Player) Volup() {
	speaker.Lock()
	p.volume.Volume += 0.1
	speaker.Unlock()
}

// Voldown decrease volume of music
func (p *Player) Voldown() {
	speaker.Lock()
	p.volume.Volume -= 0.1
	speaker.Unlock()
}

// SetVol set volume
func (p *Player) SetVol(volume float64) {
	speaker.Lock()
	p.volume.Volume = volume
	speaker.Unlock()
}

// Spdup increase speed controll
func (p *Player) Spdup() {
	speaker.Lock()
	p.resampler.SetRatio(p.resampler.Ratio() * 16 / 15)
	speaker.Unlock()
}

// Spddown decrease volume controll
func (p *Player) Spddown() {
	speaker.Lock()
	p.resampler.SetRatio(p.resampler.Ratio() * 15 / 16)
	speaker.Unlock()
}

// SetSpd set speed
func (p *Player) SetSpd(speed float64) {
	speaker.Lock()
	p.resampler.SetRatio(speed)
	speaker.Unlock()
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
