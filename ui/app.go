package ui

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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
	ppanel  *PunosPanel
	ldpanel *LoadPanel
	err     error

	// Music
	sampleRate   beep.SampleRate
	streamer     beep.StreamSeeker
	ctrl         *beep.Ctrl
	resampler    *beep.Resampler
	volume       *effects.Volume
	cuePoint     int
	musicTitle   string
	musicDirPath string

	// Waveform
	wave Waveform

	views.WidgetWatchers
}

/////////////////////////////////////////////////////
////////////////// panel transition /////////////////
/////////////////////////////////////////////////////

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

// ShowPunosPanel show PunosPanel
func (a *App) ShowPunosPanel() {
	a.show(a.ppanel)
	a.panel = a.ppanel
}

// ShowLdpanel show LoadPanel
func (a *App) ShowLdpanel() {
	a.show(a.ldpanel)
	a.panel = a.ldpanel
}

/////////////////////////////////////////////////////
////////////////// key operations ///////////////////
/////////////////////////////////////////////////////

// PlayPause is
func (a *App) PlayPause() {
	speaker.Lock()
	a.ctrl.Paused = !a.ctrl.Paused
	speaker.Unlock()
}

// Fforward fast-forward music
func (a *App) Fforward() {
	speaker.Lock()
	newPos := a.streamer.Position() + a.sampleRate.N(time.Millisecond*100)
	if newPos >= a.streamer.Len() {
		newPos = a.streamer.Len() - 1
	}
	if err := a.streamer.Seek(newPos); err != nil {
		report(err)
	}
	speaker.Unlock()
}

// Rewind rewind music
func (a *App) Rewind() {
	speaker.Lock()
	newPos := a.streamer.Position() - a.sampleRate.N(time.Millisecond*100)
	if newPos < 0 {
		newPos = 0
	}
	if err := a.streamer.Seek(newPos); err != nil {
		report(err)
	}
	speaker.Unlock()
}

// Cue set and return cue point
func (a *App) Cue() {
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

// Volup increase volume of music
func (a *App) Volup() {
	speaker.Lock()
	a.volume.Volume += 0.1
	speaker.Unlock()
}

// Voldown decrease volume of music
func (a *App) Voldown() {
	speaker.Lock()
	a.volume.Volume -= 0.1
	speaker.Unlock()
}

// SetVol set volume
func (a *App) SetVol(volume float64) {
	speaker.Lock()
	a.volume.Volume = volume
	speaker.Unlock()
}

// Spdup increase speed controll
func (a *App) Spdup() {
	speaker.Lock()
	a.resampler.SetRatio(a.resampler.Ratio() * 16 / 15)
	speaker.Unlock()
}

// Spddown decrease volume controll
func (a *App) Spddown() {
	speaker.Lock()
	a.resampler.SetRatio(a.resampler.Ratio() * 15 / 16)
	speaker.Unlock()
}

// SetSpd set speed
func (a *App) SetSpd(speed float64) {
	speaker.Lock()
	a.resampler.SetRatio(speed)
	speaker.Unlock()
}

// Status return music status
func (a *App) Status() (map[string]string, []string) {
	// gather current information
	speaker.Lock()
	pos := a.streamer.Position()
	len := a.streamer.Len()
	volume := a.volume.Volume
	speed := a.resampler.Ratio()
	speaker.Unlock()

	position := a.sampleRate.D(pos)
	length := a.sampleRate.D(len)

	cue := a.sampleRate.D(a.cuePoint)

	status := map[string]string{}

	// ToDo : building string set should move panel
	status["title"] = fmt.Sprintf("[Title : %s]", a.musicTitle)
	status["position"] = fmt.Sprintf("[Position : %s %v / %v]", GetProgressbar(a.wave.WindowSize/2, pos, len), position.Round(time.Second), length.Round(time.Second))
	status["info"] = fmt.Sprintf("[Mode\t: Normal]   [Cue Point: %v]   [Volume\t: %.1f]   [Speed\t: %.3f]", cue.Round(time.Second), volume, speed)
	status["volume"] = fmt.Sprintf("Volume\t: %.1f", volume)
	status["speed"] = fmt.Sprintf("Speed\t: %.3f", speed)
	return status, a.wave.GetWaveStr(pos)
}

// ListMusic confiel list of music
func (a *App) ListMusic() []string {
	// ToDo : use path specified by user
	cd, _ := os.Getwd()
	fileinfos, _ := ioutil.ReadDir(cd + "/mp3")
	var list []string
	r := regexp.MustCompile(`.*mp3`)
	for _, fileinfo := range fileinfos {
		if !r.MatchString(fileinfo.Name()) {
			continue
		}
		list = append(list, fileinfo.Name())
	}
	return list
}

// LoadMusic load a music to PunosPanel
func (a *App) LoadMusic(path string) {
	speaker.Lock()
	a.ctrl.Paused = true
	speaker.Unlock()
	//a.streamer.Close()

	f, err := os.Open("mp3/" + path)
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

	a.wave.Wave = LoadWave(a.wave.WaveDirPath, a.musicTitle)
	a.wave.NormalizeWave()

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

// Analyze analyze musics and create waveform
func (a *App) Analyze() {
	musicList := a.ListMusic()
	r := regexp.MustCompile(`.*mp3`)
	for _, music := range musicList {
		if !r.MatchString(music) {
			// it isn't mp3
			continue
		}
		f, err := os.Open(a.musicDirPath + "/" + music)
		if err != nil {
			report(err)
		}
		streamer, _, err := mp3.Decode(f)
		if err != nil {
			report(err)
		}
		defer streamer.Close()
		// generate and write each wave info to waveDir
		rwave := GenRawWave(streamer, a.wave.SampleInterval)
		SmoothRawWave(rwave)
		nwave := NormalizeRawWave(rwave, float64(a.wave.HeightMax), float64(a.wave.ValMax))
		WriteWave(nwave, a.wave.WaveDirPath, music)
	}
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

/////////////////////////////////////////////////////
////////////////////// mode /////////////////////////
/////////////////////////////////////////////////////

// NewApp generate new applicaiton
func NewApp() *App {
	app := &App{}
	app.app = &views.Application{}
	app.ppanel = NewPunosPanel(app)
	app.ldpanel = NewLoadPanel(app)
	app.app.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorSilver).
		Background(tcell.ColorBlack))

	// get window size
	height, width := GetWindowSize()

	//music
	//ToDo : Set Default Music
	// if there are no music in "mp3" directory,show "please insert at leaset one audio file into "mp3" directory"
	list := app.ListMusic()
	app.musicDirPath = "mp3"
	//ToDO : fix it
	if len(list) == 0 {
		fmt.Println("Please insert at lease one audio file into \"mp3\" diretctory")
		panic("error")
	}
	app.musicTitle = list[0]
	//app.musicTitle = "03.mp3"
	f, err := os.Open(app.musicDirPath + "/" + app.musicTitle)
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
	app.wave = Waveform{SampleInterval: 800, WindowSize: width, HeightMax: height / 2, ValMax: 1.0, WaveDirPath: "wave"}
	// ToDo : if there are no waveform file, call analyze function
	// ToDo : Check existance of the file reffered
	if !exists(app.wave.WaveDirPath) {
		// ToDo:create directory
		os.Mkdir("wave", 0755)
	}
	if !exists(app.wave.WaveDirPath + "/" + app.musicTitle + ".txt") {
		app.Analyze()
	}

	app.wave.Wave = LoadWave(app.wave.WaveDirPath, app.musicTitle)
	app.wave.NormalizeWave()

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

// Run the app
func (a *App) Run() {
	a.logger.Printf("Start App Running!")

	a.app.SetRootWidget(a)
	a.ShowPunosPanel()

	a.logger.Printf("Start Asynchronous function")
	// call update each second
	go func() {
		for {
			a.app.Update()
			// aim 60fps(like fighting game)
			time.Sleep(time.Millisecond * 16)
		}
	}()
	a.app.Run()
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil

}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
