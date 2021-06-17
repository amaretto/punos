package player

// ToDo: move pakage

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"regexp"
	"time"

	"github.com/amaretto/waveform/pkg/waveform"
	"github.com/benjojo/bpm"
	"github.com/dhowden/tag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/sirupsen/logrus"
)

// Analyzer is
type Analyzer struct {
	selector *Selector
	wvfmr    *waveform.Waveformer
}

func newAnalyzer(s *Selector) *Analyzer {
	a := &Analyzer{
		selector: s,
		wvfmr:    waveform.NewWaveformer(),
	}
	return a
}

// analyzeMusic analyze music reffered path
func (a *Analyzer) analyzeMusic(musicInfo *MusicInfo) {
	a.wvfmr.MusicPath = musicInfo.Path
	wvfm, err := a.wvfmr.GenWaveForm()

	// ToDo: avoid os.Exit when analyzer failed
	if err != nil {
		report(err)
	}
	wvfm.MusicTitle = musicInfo.Title

	// ToDo: analyze music info
	err = a.analyzeMusicInfo(musicInfo)
	if err != nil {
		report(err)
	}

	// set sqlite
	logrus.Debug("start create and register waveform:", musicInfo.Path)
	registerMusicInfo(musicInfo)
	registerWaveform(wvfm)
	logrus.Debug("finish create and register waveform")

	return
}

func (a *Analyzer) analyzeMusicInfo(musicInfo *MusicInfo) error {
	f, err := os.Open(musicInfo.Path)
	if err != nil {
		report(err)
	}
	defer f.Close()

	// get meta data
	m, err := tag.ReadFrom(f)
	if err != nil {
		report(err)
	}

	musicInfo.Album = m.Album()
	musicInfo.Title = m.Title()
	musicInfo.Authors = m.Artist()

	// get duration and bpm
	streamer, _, err := mp3.Decode(f)
	if err != nil {
		report(err)
	}
	defer streamer.Close()

	duration := int(a.selector.player.sampleRate.D(streamer.Len()).Round(time.Second).Seconds())
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

func (a *Analyzer) listMusic(path string) []string {
	r := regexp.MustCompile(`.*mp3`)
	cd, _ := os.Getwd()
	fileInfos, _ := os.ReadDir(cd + "/mp3")
	var list []string
	for _, fileInfo := range fileInfos {
		if !r.MatchString(fileInfo.Name()) {
			continue
		}
		list = append(list, cd+"/mp3/"+fileInfo.Name())
	}
	return list
}

func registerMusicInfo(musicInfo *MusicInfo) {
	dbPath := "mp3/test.db"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		report(err)
	}
	cmd := "INSERT INTO music(path, title, album, duration, authors, sampleRate, format, bpm) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(cmd, musicInfo.Path, musicInfo.Title, musicInfo.Album, musicInfo.Duration, musicInfo.Authors, musicInfo.SampleRate, musicInfo.Format, musicInfo.BPM)
	if err != nil {
		logrus.Debug(err)
		report(err)
	}
}

func registerWaveform(w *waveform.Waveform) {
	dbPath := "mp3/test.db"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		report(err)
	}

	data := make([]byte, len(w.Wave))
	for i, n := range w.Wave {
		data[i] = byte(n)
	}

	// ToDo: check and create databases if not exists
	cmd := "INSERT INTO waveform values(?,?)"
	_, err = db.Exec(cmd, w.MusicTitle, data)
	if err != nil {
		logrus.Debug(err)
		report(err)
	}
}
