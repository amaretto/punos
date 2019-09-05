package ui

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// App has tcell element creating ui
type App struct {
	app     *views.Application
	view    views.View
	panel   views.Widget
	logger  *log.Logger
	trntbl  *TrntblPanel
	ldpanel *LoadPanel
	err     error
	// Music
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
	cuePoint   int

	// Waveform
	waveform []int

	views.WidgetWatchers
}

// for switch Panels
func (a *App) show(w views.Widget) {
	a.app.PostFunc(func() {
		if w != a.panel {
			a.panel.SetView(nil)
			a.panel = w
		}
		a.panel.SetView(a.view)
		a.Resize()
		a.app.Refresh()
	})
}

// ShowTrntbl show trntbl Panel
func (a *App) ShowTrntbl() {
	a.Logf("in ShowTrntbl")
	a.show(a.trntbl)
	a.panel = a.trntbl
}

// ShowLdpanel show LoadPanel
func (a *App) ShowLdpanel() {
	a.Logf("in ShowLdpanel")
	a.show(a.ldpanel)
	a.panel = a.trntbl
}

/////////////////////////////////////////////////////
////////////////// key operations ///////////////////
/////////////////////////////////////////////////////

// PlayPause is
func (a *App) PlayPause() {
	a.Logf("PlayPause!!")
	speaker.Lock()
	a.ctrl.Paused = !a.ctrl.Paused
	speaker.Unlock()
}

// Fforward is fast forward module
func (a *App) Fforward() {
	speaker.Lock()
	newPos := a.streamer.Position() + a.sampleRate.N(time.Second)
	if newPos >= a.streamer.Len() {
		newPos = a.streamer.Len() - 1
	}
	if err := a.streamer.Seek(newPos); err != nil {
		report(err)
	}
	speaker.Unlock()
}

// Rewind is
func (a *App) Rewind() {
	speaker.Lock()
	newPos := a.streamer.Position() - a.sampleRate.N(time.Second)
	if newPos < 0 {
		newPos = 0
	}
	if err := a.streamer.Seek(newPos); err != nil {
		report(err)
	}
	speaker.Unlock()
}

// Cue is
func (a *App) Cue() {
	a.Logf("Cue!!")
	if a.ctrl.Paused {
		speaker.Lock()
		a.cuePoint = a.streamer.Position()
		speaker.Unlock()
	} else {
		speaker.Lock()
		a.streamer.Seek(a.cuePoint)
		speaker.Unlock()
	}
}

// Volup is volume controll
func (a *App) Volup() {
	speaker.Lock()
	a.volume.Volume += 0.1
	speaker.Unlock()
}

// Voldown is volume controll
func (a *App) Voldown() {
	speaker.Lock()
	a.volume.Volume -= 0.1
	speaker.Unlock()
}

// Spdup is speed controll
func (a *App) Spdup() {
	speaker.Lock()
	a.resampler.SetRatio(a.resampler.Ratio() * 16 / 15)
	speaker.Unlock()
}

// Spddown is volume controll
func (a *App) Spddown() {
	speaker.Lock()
	a.resampler.SetRatio(a.resampler.Ratio() * 15 / 16)
	speaker.Unlock()
}

// Status return music status
func (a *App) Status() (map[string]string, []string) {
	// gather current information
	speaker.Lock()
	pos := a.streamer.Position()
	position := a.sampleRate.D(a.streamer.Position())
	length := a.sampleRate.D(a.streamer.Len())
	volume := a.volume.Volume
	speed := a.resampler.Ratio()
	speaker.Unlock()

	cue := a.sampleRate.D(a.cuePoint)

	status := map[string]string{}
	status["position"] = fmt.Sprintf("Position\t: %v / %v", position.Round(time.Second), length.Round(time.Second))
	status["cue"] = fmt.Sprintf("Cue\t\t\t\t: %v", cue.Round(time.Second))
	status["volume"] = fmt.Sprintf("Volume\t: %.1f", volume)
	status["speed"] = fmt.Sprintf("Speed\t: %.3f", speed)
	return status, Wave2str(GetWave(a.waveform, pos, 800, 100), 30)
}

// ListMusic confiel list of music
func (a *App) ListMusic() []string {
	cd, _ := os.Getwd()
	fileinfos, _ := ioutil.ReadDir(cd + "/mp3")
	list := make([]string, len(fileinfos))
	for i, fileinfo := range fileinfos {
		// ToDo : need validation check
		list[i] = fileinfo.Name()
	}
	return list
}

// LoadMusic load a music to trntbl panel
func (a *App) LoadMusic(path string) {
	speaker.Lock()
	a.ctrl.Paused = true
	speaker.Unlock()
	//a.streamer.Close()

	f, err := os.Open("mp3/" + path)
	if err != nil {
		report(err)
	}

	//var format beep.Format
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		report(err)
	}
	// ToDo close streamer when music is switched

	// waveform
	sampleInterval := 800
	heightMax := 30
	valMax := 1.0

	wave := GenWave(streamer, sampleInterval)
	Smooth(wave)
	Smooth(wave)
	Smooth(wave)
	Smooth(wave)
	a.waveform = Normalize(wave, float64(heightMax), float64(valMax))

	a.sampleRate = format.SampleRate
	a.streamer = streamer
	a.ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, a.streamer)}
	a.resampler = beep.ResampleRatio(4, 1, a.ctrl)
	a.volume = &effects.Volume{Streamer: a.resampler, Base: 2}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))
	speaker.Play(a.volume)
	// first, pause music
	speaker.Lock()
	a.ctrl.Paused = !a.ctrl.Paused
	speaker.Unlock()
}

// Quit is
func (a *App) Quit() {
	a.app.Quit()
}

/////////////////////////////////////////////////////
////////////////////// Logs /////////////////////////
/////////////////////////////////////////////////////

// SetLogger set logger to app
func (a *App) SetLogger(logger *log.Logger) {
	a.logger = logger
	if logger != nil {
		logger.Printf("Start logger")
	}
}

// Logf print logs by referred format
func (a *App) Logf(fmt string, v ...interface{}) {
	if a.logger != nil {
		a.logger.Printf(fmt, v...)
	}
}

/////////////////////////////////////////////////////
//////////////////// handle /////////////////////////
/////////////////////////////////////////////////////

// HandleEvent handle some special key event or delegate process to Panel
func (a *App) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC:
			a.Quit()
			return true
		case tcell.KeyCtrlL:
			a.app.Refresh()
			return true
		}
	}
	if a.panel != nil {
		return a.panel.HandleEvent(ev)
	}
	return false
}

// Draw call Draw() panel has.(need it?)
func (a *App) Draw() {
	if a.panel != nil {
		a.panel.Draw()
	}
}

// Resize call Resize() panel has.(need it?)
func (a *App) Resize() {
	if a.panel != nil {
		a.panel.Resize()
	}
}

// SetView set view app have
func (a *App) SetView(view views.View) {
	a.view = view
	if a.panel != nil {
		a.panel.SetView(view)
	}
}

// Size set size of panel app have
func (a *App) Size() (int, int) {
	if a.panel != nil {
		return a.panel.Size()
	}
	return 0, 0
}

// GetAppName return application name
func (a *App) GetAppName() string {
	return "punos v0.1"
}

// NewApp generate new applicaiton
func NewApp() *App {
	app := &App{}
	app.app = &views.Application{}
	app.trntbl = NewTrntblPanel(app)
	app.ldpanel = NewLoadPanel(app)
	app.app.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorSilver).
		Background(tcell.ColorBlack))

	//music
	f, err := os.Open("mp3/01.mp3")
	if err != nil {
		report(err)
	}

	//var format beep.Format
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		report(err)
	}
	//defer streamer.Close()
	// ToDo close streamer when music is switched

	// waveform
	sampleInterval := 800
	heightMax := 30
	valMax := 1.0

	wave := GenWave(streamer, sampleInterval)
	Smooth(wave)
	Smooth(wave)
	Smooth(wave)
	Smooth(wave)
	app.waveform = Normalize(wave, float64(heightMax), float64(valMax))

	app.sampleRate = format.SampleRate
	app.streamer = streamer
	app.ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, app.streamer)}
	app.resampler = beep.ResampleRatio(4, 1, app.ctrl)
	app.volume = &effects.Volume{Streamer: app.resampler, Base: 2}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))
	speaker.Play(app.volume)
	// first, pause music
	speaker.Lock()
	app.ctrl.Paused = !app.ctrl.Paused
	speaker.Unlock()

	//go app.refresh()
	return app
}

//func (a *App)refresh() {
//
//}

// Run the app
func (a *App) Run() {
	a.Logf("Punos")
	a.app.SetRootWidget(a)
	a.ShowTrntbl()
	// call update each second
	go func() {
		for {
			a.app.Update()
			// if set time.Millisecond, this app freez...
			time.Sleep(time.Millisecond * 2)
		}
	}()
	a.Logf("Starting app loop")
	a.app.Run()
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
