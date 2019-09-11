package ui

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/garyburd/redigo/redis"
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
	sampleRate   beep.SampleRate
	streamer     beep.StreamSeeker
	ctrl         *beep.Ctrl
	resampler    *beep.Resampler
	volume       *effects.Volume
	cuePoint     int
	musicTitle   string
	musicDirPath string

	// Waveform
	waveform       []int
	sampleInterval int
	windowSize     int
	heightMax      int
	valMax         float64
	waveDirPath    string

	// Mode
	Mode         string
	b2bAvailable bool
	// b2b mode
	b2bTarget string // redis ip and port
	con       redis.Conn
	// sync Mode
	isSync     bool
	syncTarget string

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
	if a.Mode == "sync" {
		//setPos
		pos := a.streamer.Position()
		posstr := strconv.Itoa(pos)
		redisSet("pos", posstr, a.con)
		//set play/pause
		if a.ctrl.Paused {
			redisSet("play", "0", a.con)
		} else {
			redisSet("play", "1", a.con)
		}
	}
	speaker.Unlock()
}

// Fforward is fast forward module
func (a *App) Fforward() {
	speaker.Lock()
	newPos := a.streamer.Position() + a.sampleRate.N(time.Millisecond*100)
	if newPos >= a.streamer.Len() {
		newPos = a.streamer.Len() - 1
	}
	if err := a.streamer.Seek(newPos); err != nil {
		report(err)
	}
	if a.Mode == "sync" {
		//setPos
		pos := a.streamer.Position()
		posstr := strconv.Itoa(pos)
		redisSet("pos", posstr, a.con)
	}
	speaker.Unlock()
}

// Rewind is
func (a *App) Rewind() {
	speaker.Lock()
	newPos := a.streamer.Position() - a.sampleRate.N(time.Millisecond*100)
	if newPos < 0 {
		newPos = 0
	}
	if err := a.streamer.Seek(newPos); err != nil {
		report(err)
	}
	if a.Mode == "sync" {
		//setPos
		pos := a.streamer.Position()
		posstr := strconv.Itoa(pos)
		redisSet("pos", posstr, a.con)
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
		if a.Mode == "sync" {
			pos := a.streamer.Position()
			posstr := strconv.Itoa(pos)
			redisSet("pos", posstr, a.con)
		}
		speaker.Unlock()
	}
}

// Volup is volume controll
func (a *App) Volup() {
	speaker.Lock()
	a.volume.Volume += 0.1
	if a.Mode == "sync" {
		volume := strconv.FormatFloat(a.volume.Volume, 'f', 4, 64)
		redisSet("volume", volume, a.con)
	}
	speaker.Unlock()
}

// Voldown is volume controll
func (a *App) Voldown() {
	speaker.Lock()
	a.volume.Volume -= 0.1
	if a.Mode == "sync" {
		volume := strconv.FormatFloat(a.volume.Volume, 'f', 4, 64)
		redisSet("volume", volume, a.con)
	}
	speaker.Unlock()
}

// SetVol is volume controll
func (a *App) SetVol(volume float64) {
	speaker.Lock()
	a.volume.Volume = volume
	speaker.Unlock()
}

// Spdup is speed controll
func (a *App) Spdup() {
	speaker.Lock()
	a.resampler.SetRatio(a.resampler.Ratio() * 16 / 15)
	if a.Mode == "sync" {
		speed := strconv.FormatFloat(a.resampler.Ratio(), 'f', 4, 64)
		redisSet("speed", speed, a.con)
	}
	speaker.Unlock()
}

// Spddown is volume controll
func (a *App) Spddown() {
	speaker.Lock()
	a.resampler.SetRatio(a.resampler.Ratio() * 15 / 16)
	if a.Mode == "sync" {
		speed := strconv.FormatFloat(a.resampler.Ratio(), 'f', 4, 64)
		redisSet("speed", speed, a.con)
	}
	speaker.Unlock()
}

// SetSpd is volume controll
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

	status["title"] = fmt.Sprintf("[Title : %s]", a.musicTitle)
	status["position"] = fmt.Sprintf("[Position : %s %v / %v]", GetProgressbar(a.windowSize/2, pos, len), position.Round(time.Second), length.Round(time.Second))
	status["info"] = fmt.Sprintf("[Mode\t: %s]   [Cue Point: %v]   [Volume\t: %.1f]   [Speed\t: %.3f]", a.Mode, cue.Round(time.Second), volume, speed)
	status["volume"] = fmt.Sprintf("Volume\t: %.1f", volume)
	status["speed"] = fmt.Sprintf("Speed\t: %.3f", speed)
	return status, Wave2str(GetWave(a.waveform, pos, a.sampleInterval, a.windowSize), a.heightMax)
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
	// update title
	a.musicTitle = path

	//var format beep.Format
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		report(err)
	}
	// ToDo close streamer when music is switched

	//wave := GenWave(streamer, a.sampleInterval)
	//Smooth(wave)
	//a.waveform = Normalize(wave, float64(a.heightMax), float64(a.valMax))
	a.waveform = LoadWave(a.waveDirPath, a.musicTitle)

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
	if a.Mode == "sync" {
		redisSet("title", a.musicTitle, a.con)
	}
	speaker.Unlock()
}

// Analyze is
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
		wave := GenWave(streamer, a.sampleInterval)
		Smooth(wave)
		nwave := Normalize(wave, float64(a.heightMax), float64(a.valMax))
		WriteWave(nwave, a.waveDirPath, music)
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

// SetMode is
func (a *App) SetMode(mode string) {
	if !a.b2bAvailable {
		a.Mode = "normal"
	} else {
		a.Mode = mode
	}
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
	app.musicDirPath = "mp3"
	app.musicTitle = "03.mp3"
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
	app.waveDirPath = "wave"
	app.windowSize = 141
	app.sampleInterval = 800
	app.heightMax = 25
	app.valMax = 1.0
	app.waveform = LoadWave(app.waveDirPath, app.musicTitle)

	// mode
	app.Mode = "normal"
	app.b2bTarget = "192.168.1.40:6379"

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

type winsize struct {
	Row    int
	Col    int
	Xpixel int
	Ypixel int
}

func getWidth() int {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return int(ws.Col)
}

// Run the app
func (a *App) Run() {
	a.Logf("Punos")
	a.app.SetRootWidget(a)
	a.ShowTrntbl()

	// for b2b mode
	// need to see difference
	var tmpPlay string
	var newPlay string
	var tmpTitle string
	var newTitle string
	var tmppos string
	var newpos string
	var tmpVolume string
	var newVolume string
	var tmpSpeed string
	var newSpeed string

	var err error
	a.con, err = redisConnection(a.b2bTarget)
	// if redis isn't exist, b2b is not available
	if err != nil {
		a.b2bAvailable = false
	} else {
		a.b2bAvailable = true
		//isPaused(0:pause,1:play)
		tmpPlay = "0"
		newPlay = "0"
		redisSet("play", "0", a.con)
		//title
		tmpTitle = a.musicTitle
		redisSet("title", a.musicTitle, a.con)
		//pos
		tmppos = "0"
		redisSet("pos", "0", a.con)
		//volume
		tmpVolume := strconv.FormatFloat(a.volume.Volume, 'f', 4, 64)
		redisSet("volume", tmpVolume, a.con)
		//speed
		tmpSpeed := strconv.FormatFloat(a.resampler.Ratio(), 'f', 4, 64)
		redisSet("speed", tmpSpeed, a.con)
	}

	// call update each second
	go func() {
		for {
			a.app.Update()
			// aim 60fps(like fighting game)
			time.Sleep(time.Millisecond * 16)
			// for b2b mode
			if a.Mode == "b2b" {
				//isPaused
				newPlay = redisGet("play", a.con)
				if newPlay != tmpPlay {
					a.PlayPause()
					tmpPlay = newPlay
					redisSet("play", tmpPlay, a.con)
				}
				//title
				newTitle = redisGet("title", a.con)
				if newTitle != tmpTitle {
					a.LoadMusic(newTitle)
					tmpTitle = newTitle
					redisSet("title", tmpTitle, a.con)
				}
				//position
				newpos = redisGet("pos", a.con)
				if newpos != tmppos {
					pos, _ := strconv.Atoi(newpos)
					a.streamer.Seek(pos)
					// reset position
					tmppos = "0"
					newpos = "0"
					redisSet("pos", "0", a.con)
				}
				//volume
				newVolume = redisGet("volume", a.con)
				if newVolume != tmpVolume {
					volume, _ := strconv.ParseFloat(newVolume, 64)
					a.SetVol(volume)
					tmpVolume = newVolume
					redisSet("volume", tmpVolume, a.con)
				}
				//speed
				newSpeed = redisGet("speed", a.con)
				if newSpeed != tmpSpeed {
					speed, _ := strconv.ParseFloat(newSpeed, 64)
					a.SetSpd(speed)
					tmpSpeed = newSpeed
					redisSet("speed", tmpSpeed, a.con)
				}
			}
		}
	}()
	a.Logf("Starting app loop")
	a.app.Run()
}

// Redis functions
func redisConnection(target string) (redis.Conn, error) {

	c, err := redis.Dial("tcp", target)
	return c, err
}

func redisSet(key string, value string, c redis.Conn) {
	c.Do("SET", key, value)
}

func redisGet(key string, c redis.Conn) string {
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return s
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
