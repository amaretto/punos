package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"
)

func drawTextLine(screen tcell.Screen, x, y int, s string, style tcell.Style) {
	for _, r := range s {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

type punosPanel struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
}

func newPunosPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker) *punosPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &punosPanel{sampleRate, streamer, ctrl, resampler, volume}
}

func (pp *punosPanel) play() {
	speaker.Play(pp.volume)
}

func (pp *punosPanel) draw(screen tcell.Screen) {
	mainStyle := tcell.StyleDefault.
		Background(tcell.NewHexColor(0x473437)).
		Foreground(tcell.NewHexColor(0xD7D8A2))
	statusStyle := mainStyle.
		Foreground(tcell.NewHexColor(0xDDC074)).
		Bold(true)

	screen.Fill(' ', mainStyle)

	drawTextLine(screen, 0, 0, " _ __  _   _ _ __   ___  ___", mainStyle)
	drawTextLine(screen, 0, 1, "| '_ \\| | | | '_ \\ / _ \\/ __|", mainStyle)
	drawTextLine(screen, 0, 2, "| |_) | |_| | | | | (_) \\__ \\", mainStyle)
	drawTextLine(screen, 0, 3, "| .__/ \\__,_|_| |_|\\___/|___/", mainStyle)
	drawTextLine(screen, 0, 4, "|_|", mainStyle)

	drawTextLine(screen, 0, 6, "Press [ESC] to quit", mainStyle)
	drawTextLine(screen, 0, 7, "Press [SPACE] to pause/resume", mainStyle)
	drawTextLine(screen, 0, 8, "Use keys in (?/?) to turn the buttons.", mainStyle)

	speaker.Lock()
	position := pp.sampleRate.D(pp.streamer.Position())
	length := pp.sampleRate.D(pp.streamer.Len())
	volume := pp.volume.Volume
	speed := pp.resampler.Ratio()
	speaker.Unlock()

	positionStatus := fmt.Sprintf("%v / %v", position.Round(time.Second), length.Round(time.Second))
	volumeStatus := fmt.Sprintf("%.1f", volume)
	speedStatus := fmt.Sprintf("%.3fx", speed)

	drawTextLine(screen, 0, 10, "Position(Q/W):", mainStyle)
	drawTextLine(screen, 16, 10, positionStatus, statusStyle)

	drawTextLine(screen, 0, 11, "Volume (A/S):", mainStyle)
	drawTextLine(screen, 16, 11, volumeStatus, statusStyle)

	drawTextLine(screen, 0, 12, "Speed (Z/X):", mainStyle)
	drawTextLine(screen, 16, 12, speedStatus, statusStyle)
}

func (pp *punosPanel) handle(event tcell.Event) (changed, quit bool) {
	switch event := event.(type) {
	case *tcell.EventKey:
		if event.Key() == tcell.KeyESC {
			return false, true
		}

		if event.Key() != tcell.KeyRune {
			return false, false
		}

		switch unicode.ToLower(event.Rune()) {
		case ' ':
			speaker.Lock()
			pp.ctrl.Paused = !pp.ctrl.Paused
			speaker.Unlock()
			return false, false

		// seek
		case 'q', 'w':
			speaker.Lock()
			newPos := pp.streamer.Position()
			if event.Rune() == 'q' {
				newPos -= pp.sampleRate.N(time.Second)
			}
			if event.Rune() == 'w' {
				newPos += pp.sampleRate.N(time.Second)
			}
			if newPos >= pp.streamer.Len() {
				newPos = pp.streamer.Len() - 1
			}
			if err := pp.streamer.Seek(newPos); err != nil {
				report(err)
			}
			speaker.Unlock()
			return true, false

		// volume
		case 'a':
			speaker.Lock()
			pp.volume.Volume -= 0.1
			speaker.Unlock()
			return true, false

		case 's':
			speaker.Lock()
			pp.volume.Volume += 0.1
			speaker.Unlock()
			return true, false

		// change speed
		case 'z':
			speaker.Lock()
			pp.resampler.SetRatio(pp.resampler.Ratio() * 15 / 16)
			speaker.Unlock()
			return true, false

		case 'x':
			speaker.Lock()
			pp.resampler.SetRatio(pp.resampler.Ratio() * 16 / 15)
			speaker.Unlock()
			return true, false
		}
	}
	return false, false
}

func main() {

	f, err := os.Open("mp3/01.mp3")
	if err != nil {
		log.Fatal(err)
	}

	//2nd
	f2, err := os.Open("mp3/02.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	//2nd
	streamer2, format, err := mp3.Decode(f2)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer2.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	screen, err := tcell.NewScreen()
	if err != nil {
		report(err)
	}
	err = screen.Init()
	if err != nil {
		report(err)
	}
	defer screen.Fini()

	pp := newPunosPanel(format.SampleRate, streamer)
	//2nd
	pp2 := newPunosPanel(format.SampleRate, streamer2)

	screen.Clear()
	pp.draw(screen)
	screen.Show()

	pp.play()
	pp2.play()

	seconds := time.Tick(time.Second)
	events := make(chan tcell.Event)
	go func() {
		for {
			events <- screen.PollEvent()
		}
	}()

loop:
	for {
		select {
		case event := <-events:
			changed, quit := pp.handle(event)
			if quit {
				break loop
			}
			if changed {
				screen.Clear()
				pp.draw(screen)
				screen.Show()
			}
		case <-seconds:
			screen.Clear()
			pp.draw(screen)
			screen.Show()
		}
	}
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
