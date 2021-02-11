// ToDo: move pakage
package player

import "github.com/amaretto/waveform/pkg/waveform"

type Analyzer struct {
	wfmr *waveform.Waveformer
}

type MusicInfo struct {
	wave []int
}

func newAnalyzer() *Analyzer {
	a := &Analyzer{
		wfmr: waveform.NewWaveformer(),
	}
	return a
}

// analyzeDir anlyze musics in reffered directory and create info record to sqlite
func (a *Analyzer) analyzeDir(path string) {
	// ToDo : list musics
	// ToDo : call analyzeMusic()

}

// analyzeMusic analyze music reffered path
func (a *Analyzer) analyzeMusic(path string) {
	// ToDo : load music
	// ToDo :
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
