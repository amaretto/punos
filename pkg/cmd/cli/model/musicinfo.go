package model

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

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
	List []*MusicInfo
	DB   *sql.DB
	conf *config.Config
}

func NewMusics(conf *config.Config) *Musics {
	db, err := sql.Open("sqlite3", conf.DBPath)
	if err != nil {
		report(err)
	}
	//ToDo: fix it?
	//defer db.Close()
	musics := &Musics{conf: conf, DB: db}
	return musics
}

func (m *Musics) Load() {
	rows, err := m.DB.Query("SELECT path, title, album, authors, duration, sampleRate, format, bpm FROM music")
	if err != nil {
		report(err)
	}

	m.List = []*MusicInfo{}
	// retrieve data from DB
	for rows.Next() {
		mi := &MusicInfo{}
		if err := rows.Scan(&mi.Path, &mi.Title, &mi.Album, &mi.Authors, &mi.Duration, &mi.SampleRate, &mi.Format, &mi.BPM); err != nil {
			report(err)
		}
		m.List = append(m.List, mi)
	}

	// check music status
	musicPathList := listMusic(m.conf.MusicPath)
	m.checkMusicStatus(musicPathList)
}

func (m *Musics) checkMusicStatus(musicPathList []string) {
	// data exists && no file
	for _, mi := range m.List {
		if contains(mi.Path, musicPathList) {
			mi.Status = "✔"
			musicPathList = del(mi.Path, musicPathList)
		} else {
			mi.Status = "Moved"
		}
	}

	// file exists && no data
	for _, musicPath := range musicPathList {
		mi := &MusicInfo{}
		mi.Path = musicPath
		mi.Title = filepath.Base(musicPath)
		mi.Status = "Not Analyzed"
		m.List = append(m.List, mi)
	}
}

func contains(path string, list []string) bool {
	for _, s := range list {
		if s == path {
			return true
		}
	}
	return false
}

func listMusic(musicPath string) []string {
	r := regexp.MustCompile(`.*mp3`)
	fileInfos, _ := os.ReadDir(musicPath)
	var list []string
	for _, fileInfo := range fileInfos {
		if !r.MatchString(fileInfo.Name()) {
			continue
		}
		list = append(list, musicPath+"/"+fileInfo.Name())
	}
	return list
}

func del(path string, list []string) []string {
	for i, s := range list {
		if s == path {
			if i < len(list)-1 {
				return append(list[:i], list[i+1:]...)
			}
			return list[:i]
		}
	}
	return list
}

func (m Musics) LoadWaveform(mi *MusicInfo) {
	cmd := "SELECT wave FROM waveform WHERE path = ?"

	row := m.DB.QueryRow(cmd, filepath.Base(mi.Path))
	var data []byte
	row.Scan(&data)
	mi.Waveform = data
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
