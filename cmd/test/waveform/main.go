package main

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	gomp3 "github.com/hajimehoshi/go-mp3"
)

const (
	gomp3NumChannels   = 2
	gomp3Precision     = 2
	gomp3BytesPerFrame = gomp3NumChannels * gomp3Precision
)

func main() {

	f, err := os.Open("mp3/02.mp3")
	if err != nil {
		report(err)
	}
	d, err := gomp3.NewDecoder(f)
	if err != nil {
		report(err)
	}

	format := beep.Format{
		SampleRate:  beep.SampleRate(d.SampleRate()),
		NumChannels: gomp3NumChannels,
		Precision:   gomp3Precision,
	}

	bufferSize := format.SampleRate.N(time.Second / 30)
	fmt.Println(bufferSize)

	//samples := make([][2]float64, bufferSize)

}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
