package player

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const logo = `





________   ___  ___   ________    ________   ________      
|\   __ \ |\  \|\  \ |\   ___  \ |\   __  \ |\   ____\     
\ \  \|\ \ \ \  \\  \\ \  \\ \  \\ \  \|\  \\ \  \___|_    
 \ \   ___\ \ \  \\  \\ \  \\ \  \\ \  \\\  \\ \_____  \   
  \ \  \___| \ \  \\  \\ \  \\ \  \\ \  \\\  \\|____|\  \  
   \ \__\     \ \_______\\\__\\ \__\\ \_______\ ____\_\  \ 
    \|__|      \|_______| \|__| \|__| \|_______||\________\
                                                \|________|
`

type WaveformPanel struct {
	*DefaultView
}

func NewWaveformPanel() *WaveformPanel {
	w := &WaveformPanel{
		DefaultView: NewDefaultView("Waveform"),
	}
	w.initWaveformPanel()
	return w
}

func (w *WaveformPanel) initWaveformPanel() {
	w.SetText(logo)
}

func (w *WaveformPanel) loadWaveform(title string) {
	dbPath := "mp3/test.db"

	// ToDo: Implement error handling
	con, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalln(err)
	}

	cmd := "SELECT wave FROM waveform WHERE title = ?"

	_, err = con.Exec(cmd, title)
	if err != nil {
		log.Fatalln(err)
	}

}
