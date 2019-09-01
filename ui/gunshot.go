package ui

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// Gunshot is
type Gunshot struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
	buffer     *beep.Buffer
}

// Init is
func (gs *Gunshot) Init() {
	f, err := os.Open("gunshot.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	gs.sampleRate = format.SampleRate
	gs.streamer = streamer
	gs.ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, gs.streamer)}
	gs.resampler = beep.ResampleRatio(4, 1, gs.ctrl)
	gs.volume = &effects.Volume{Streamer: gs.resampler, Base: 2}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/60))

	gs.buffer = beep.NewBuffer(format)
	gs.buffer.Append(streamer)
	streamer.Close()
}

// Shot is
func (gs *Gunshot) Shot() {
	for {
		fmt.Print("Press [ENTER] to fire a gunshot! ")
		fmt.Scanln()

		shot := gs.buffer.Streamer(0, gs.buffer.Len())
		speaker.Play(shot)
	}
}

// NewGunshot is
func NewGunshot() *Gunshot {
	gs := &Gunshot{}
	gs.Init()
	return gs
}
