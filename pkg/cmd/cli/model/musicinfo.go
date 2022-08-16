package model

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/amaretto/punos/pkg/cmd/cli/config"
	"github.com/amaretto/waveform/pkg/waveform"
	_ "github.com/mattn/go-sqlite3"
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
	List       []*MusicInfo
	db         *sql.DB
	conf       *config.Config
	currentDir string
}

func NewMusics(conf *config.Config) *Musics {
	db, err := sql.Open("sqlite3", conf.DBPath)
	if err != nil {
		report(err)
	}

	musics := &Musics{currentDir: conf.MusicPath, db: db}
	musics.initDB()
	return musics
}

func (m *Musics) initDB() {
	//music
	_, err := m.db.Exec("CREATE TABLE IF NOT EXISTS music(path text PRIMARY KEY, title text, album text, duration text, authors text, sampleRate integer, format text, bpm[interger] DEFAULT 100)")
	if err != nil {
		report(err)
	}
	//waveform
	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS waveform(path text PRIMARY KEY, wave blob)")
	if err != nil {
		report(err)
	}

}

// ToDo: implement functions allow to move between directories
func (m *Musics) ListMusics() {
	musicPathList := listMusicPath(m.currentDir)

	// retrieve data from DB
	// ToDo: filter
	rows, err := m.db.Query("SELECT path, title, album, authors, duration, sampleRate, format, bpm FROM music")
	if err != nil {
		report(err)
	}

	m.List = []*MusicInfo{}
	for rows.Next() {
		mi := &MusicInfo{}
		if err := rows.Scan(&mi.Path, &mi.Title, &mi.Album, &mi.Authors, &mi.Duration, &mi.SampleRate, &mi.Format, &mi.BPM); err != nil {
			report(err)
		}
		if contains(mi.Path, musicPathList) {
			mi.Status = "Ready"
			m.List = append(m.List, mi)
			musicPathList = del(mi.Path, musicPathList)
		}
	}
	m.appendNotAnalyzedMusics(musicPathList)
}

func (m *Musics) RegisterMusicInfo(musicInfo *MusicInfo) {
	cmd := "INSERT INTO music(path, title, album, duration, authors, sampleRate, format, bpm) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := m.db.Exec(cmd, musicInfo.Path, musicInfo.Title, musicInfo.Album, musicInfo.Duration, musicInfo.Authors, musicInfo.SampleRate, musicInfo.Format, musicInfo.BPM)
	if err != nil {
		report(err)
	}
}

func (m *Musics) RegisterWaveform(w *waveform.Waveform) {
	data := make([]byte, len(w.Wave))
	for i, n := range w.Wave {
		data[i] = byte(n)
	}
	cmd := "INSERT INTO waveform values(?,?)"
	_, err := m.db.Exec(cmd, w.MusicTitle, data)
	if err != nil {
		report(err)
	}
}

func (m *Musics) appendNotAnalyzedMusics(musicPathList []string) {
	for _, musicPath := range musicPathList {
		mi := &MusicInfo{}
		mi.Path = musicPath
		mi.Title = filepath.Base(musicPath)
		mi.Status = "Not Analyzed"
		m.List = append(m.List, mi)
	}
}

// Musics Dirs
func (m *Musics) ListDirs() []string {
	fileInfos, _ := os.ReadDir(m.currentDir)
	var list []string
	for _, fileInfo := range fileInfos {
		if fileInfo.Type().IsDir() {
			list = append(list, fileInfo.Name())
		}
	}
	return list
}

func (m *Musics) MoveChildDir(dirName string) {
	m.currentDir = m.currentDir + "/" + dirName
	// update music info
	m.ListMusics()
}

func (m *Musics) MoveParentDir() {
	if m.currentDir != "/" {
		m.currentDir = filepath.Dir(m.currentDir)
	}
	m.ListMusics()
}

func contains(path string, list []string) bool {
	for _, s := range list {
		if s == path {
			return true
		}
	}
	return false
}

func listMusicPath(musicPath string) []string {
	r := regexp.MustCompile(`.*mp3`)
	fileInfos, _ := os.ReadDir(musicPath)
	var list []string
	for _, fileInfo := range fileInfos {
		if !r.MatchString(fileInfo.Name()) || fileInfo.Type().IsDir() {
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

	row := m.db.QueryRow(cmd, mi.Path)
	var data []byte
	row.Scan(&data)
	mi.Waveform = data
}

const SampleInterval = 800

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
