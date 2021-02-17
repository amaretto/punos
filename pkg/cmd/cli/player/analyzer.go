// ToDo: move pakage
package player

import (
	"database/sql"
	"os"
	"regexp"

	"github.com/amaretto/waveform/pkg/waveform"
)

type Analyzer struct {
	wvfmr *waveform.Waveformer
}

type MusicInfo struct {
	wave []int
}

func newAnalyzer() *Analyzer {
	a := &Analyzer{
		wvfmr: waveform.NewWaveformer(),
	}
	return a
}

// analyzeDir anlyze musics in reffered directory and create info record to sqlite
func (a *Analyzer) analyzeDir(path string) {
	// ToDo: fix
	musicList := a.listMusic("dummy")
	finished := make(chan bool)
	for _, music := range musicList {
		musicPath := music
		go func() {
			a.analyzeMusic(musicPath)
			finished <- true
		}()
	}
	for i := 0; i < len(musicList); i++ {
		<-finished
	}
}

// analyzeMusic analyze music reffered path
func (a *Analyzer) analyzeMusic(musicPath string) {
	// generate and write each wave info to waveDir
	a.wvfmr.MusicPath = musicPath
	wvfm, err := a.wvfmr.GenWaveForm()
	// ToDo: avoid os.Exit when analyzer failed
	if err != nil {
		report(err)
	}
	// ToDo: set sqlite
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

func registerWaveform(w *waveform.Waveform) {
	dbPath := "mp3/test.db"

	con, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		report(err)
	}

	// ToDo: create table if not exist

	cmd := "INSERT INTO waveform values('?',json_array(?))"

	// ToDo: convert Wave to string
	_, err = con.Exec(cmd, w.MusicTitle, w.Wave)
	if err != nil {
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
//
//func report(err error) {
//	fmt.Fprintln(os.Stderr, err)
//	os.Exit(1)
//}
