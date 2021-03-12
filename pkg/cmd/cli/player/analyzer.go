package player

// ToDo: move pakage

import (
	"database/sql"
	"os"
	"regexp"

	"github.com/amaretto/waveform/pkg/waveform"
	"github.com/sirupsen/logrus"
)

// Analyzer is
type Analyzer struct {
	wvfmr *waveform.Waveformer
}

func newAnalyzer() *Analyzer {
	a := &Analyzer{
		wvfmr: waveform.NewWaveformer(),
	}
	return a
}

// analyzeDir anlyze musics in reffered directory and create info record to sqlite
//func (a *Analyzer) analyzeDir(path string) {
//	// ToDo: fix
//	musicList := a.listMusic("dummy")
//	finished := make(chan bool)
//	for _, music := range musicList {
//		musicPath := music
//		go func() {
//			a.analyzeMusic(musicPath)
//			finished <- true
//		}()
//	}
//	for i := 0; i < len(musicList); i++ {
//		<-finished
//	}
//}

// analyzeMusic analyze music reffered path
func (a *Analyzer) analyzeMusic(musicInfo *MusicInfo) {
	// generate and write each wave info to waveDir
	a.wvfmr.MusicPath = musicInfo.Path
	wvfm, err := a.wvfmr.GenWaveForm()
	// ToDo: avoid os.Exit when analyzer failed
	if err != nil {
		report(err)
	}
	wvfm.MusicTitle = musicInfo.Title

	// set sqlite
	logrus.Debug("registerrrrrrrrrr")
	registerMusicInfo(musicInfo)
	registerWaveform(wvfm)
	return
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
	cmd := "INSERT INTO music(path, title, album, duration, authors, sampleRate, format) VALUES(?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(cmd, musicInfo.Path, musicInfo.Title, musicInfo.Album, musicInfo.Duration, musicInfo.Authors, musicInfo.SampleRate, musicInfo.Format)
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

//func loadWaveform(title string) {
//	dbPath := "mp3/test.db"
//
//	// ToDo: Implement error handling
//	con, err := sql.Open("sqlite3", dbPath)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	cmd := "SELECT wave FROM waveform WHERE title = ?"
//
//	_, err = con.Exec(cmd, title)
//	if err != nil {
//		log.Fatalln(err)
//	}
//}
