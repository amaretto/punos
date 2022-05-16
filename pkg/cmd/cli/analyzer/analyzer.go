package analyzer

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/amaretto/punos/pkg/cmd/cli/config"
	"github.com/amaretto/punos/pkg/cmd/cli/model"
	"github.com/amaretto/waveform/pkg/waveform"
	"github.com/benjojo/bpm"
	"github.com/dhowden/tag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/sirupsen/logrus"
)

// Analyzer is
type Analyzer struct {
	sampleRate beep.SampleRate
	wvfmr      *waveform.Waveformer
	musics     *model.Musics
}

func NewAnalyzer(conf *config.Config, sampleRate beep.SampleRate) *Analyzer {
	a := &Analyzer{
		sampleRate: sampleRate,
		wvfmr:      waveform.NewWaveformer(),
		musics:     model.NewMusics(conf),
	}
	return a
}

// analyzeMusic analyze music reffered path
func (a *Analyzer) AnalyzeMusic(musicInfo *model.MusicInfo) {
	a.wvfmr.MusicPath = musicInfo.Path
	wvfm, err := a.wvfmr.GenWaveForm()

	// ToDo: avoid os.Exit when analyzer failed
	if err != nil {
		logrus.Debug(err)
		report(err)
	}
	wvfm.MusicTitle = musicInfo.Path

	err = a.analyzeMusicInfo(musicInfo)
	if err != nil {
		logrus.Debug(err)
		report(err)
	}

	// set sqlite
	logrus.Debug("start create and register waveform:", musicInfo.Path)
	a.musics.RegisterMusicInfo(musicInfo)
	a.musics.RegisterWaveform(wvfm)
	logrus.Debug("finish create and register waveform")

	return
}

func (a *Analyzer) analyzeMusicInfo(musicInfo *model.MusicInfo) error {
	logrus.Debug("analyze")
	f, err := os.Open(musicInfo.Path)
	if err != nil {
		logrus.Debug(err)
		report(err)
	}
	defer f.Close()

	// get meta data
	logrus.Debug("get meta")
	m, err := tag.ReadFrom(f)
	if err != nil {
		logrus.Debug(err)
		report(err)
	}

	musicInfo.Album = m.Album()
	musicInfo.Title = m.Title()
	musicInfo.Authors = m.Artist()

	// get duration and bpm
	streamer, _, err := mp3.Decode(f)
	if err != nil {
		logrus.Debug(err)
		report(err)
	}
	defer streamer.Close()

	logrus.Debug("get bpm")
	duration := int(a.sampleRate.D(streamer.Len()).Round(time.Second).Seconds())
	musicInfo.Duration = fmt.Sprintf("%d:%02d", duration/60, duration%60)

	musicInfo.BPM = fmt.Sprintf("%.2f", a.detectBPM(streamer))
	logrus.Debug(musicInfo.Title, ":", musicInfo.BPM)

	return nil
}

func calcChunkLen(second int) int {
	return (bpm.RATE / bpm.INTERVAL) * second
}

func (a *Analyzer) detectBPM(streamer beep.StreamSeeker) float64 {
	var tmp [2][2]float64
	var samplel, sampler float64
	var n, v, z float64
	nrg := make([]float32, 0)
	bpms := make([]float64, 0)
	maxsize := calcChunkLen(10)
	for {
		if sn, sok := streamer.Stream(tmp[:1]); sn == 0 && !sok {
			break
		}
		samplel = tmp[0][0]
		sampler = tmp[0][1]
		z = (math.Abs(samplel) + math.Abs(sampler)) / 2
		if z > v {
			v += (z - v) / 8
		} else {
			v -= (v - z) / 512
		}

		n++
		if n == bpm.INTERVAL {
			n = 0
			nrg = append(nrg, float32(v))
		}

		if len(nrg) == maxsize {
			bpms = append(bpms, bpm.ScanForBpm(nrg, 90, 200, 1024, 1024))
			nrg = make([]float32, 0)
		}
	}
	// ToDo:sort
	return bpms[len(bpms)/2]
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
