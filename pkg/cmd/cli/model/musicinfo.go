package model

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/amaretto/punos/pkg/cmd/cli/config"
)

// MusicInfo have details of mp3 files
type MusicInfo struct {
	Status     string
	Path       string
	Album      string
	Title      string
	Authors    string
	Duration   string
	SampleRate int
	Format     string
	Waveform   []byte
	BPM        string
}

type Musics struct {
	list []*MusicInfo
	db   *sql.DB
}

func New(conf config.Config) *Musics {
	musics := &Musics{}
	db, err := sql.Open("sqlite3", conf.DBPath)
	if err != nil {
		report(err)
	}
	defer db.Close()
	musics.db = db
	return musics
}

// ToDo: Implement
func (m *Musics) GetMusics(path string) *Musics {
	return nil
}

// ToDo: Implement
func (m *Musics) loadWaveform(path string) {

}

const SampleInterval = 800

// ToDo: Implement
func initDB(confPath string) error {
	return nil
}

func (m *MusicInfo) getWaveStr(pos, width, height int) []string {
	return m.wave2str(m.getWave(pos, SampleInterval, width), height)
}

// GetWave return part of wave by position
func (m *MusicInfo) getWave(pos, sampleInterval, windowSize int) []byte {
	var wave []byte
	wave = make([]byte, windowSize)
	center := int(pos / sampleInterval)

	var idx int
	for i := 0; i < windowSize; i++ {
		idx = center - (windowSize/2 - 1) + i
		// front
		if idx < 0 {
			continue
		}
		// back
		if idx >= len(m.Waveform) {
			break
		}
		wave[i] = m.Waveform[idx]
	}
	return wave
}

// wave2str generate strings express waveform
func (m *MusicInfo) wave2str(wave []byte, limit int) []string {
	var waveStr []string
	var fill []bool
	waveStr = make([]string, limit)
	fill = make([]bool, len(wave))

	for i := limit; i > 0; i-- {
		str := ""
		for j, num := range wave {
			if j == len(wave)/2-1 {
				str = str + "|"
			} else if int(num) >= i || fill[num] {
				str = str + "#"
				fill[num] = true
			} else {
				str = str + " "
			}
		}
		waveStr[limit-i] = str
	}
	return waveStr
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
