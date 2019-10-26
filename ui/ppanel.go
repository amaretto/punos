package ui

import (
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var (
	// StyleNormal is
	StyleNormal = tcell.StyleDefault.
		Foreground(tcell.NewHexColor(0x00FF2B)).
		//			Foreground(tcell.ColorSilver).
		Background(tcell.ColorBlack)
)

// PunosPanel is
type PunosPanel struct {
	width   int
	height  int
	curx    int
	cury    int
	content *views.CellView
	text    *views.TextArea
	lines   []string
	styles  []tcell.Style
	Panel
}

// HandleEvent is
func (p *PunosPanel) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEsc:
			p.App().Quit()
			return true
		}

		switch unicode.ToLower(ev.Rune()) {
		case ' ':
			p.App().PlayPause()
			return true
		case 'c':
			p.App().Cue()
			return true
		case 'w':
			p.App().Fforward()
			return true
		case 'q':
			p.App().Rewind()
			return true
		// volume
		case 'a':
			p.App().Voldown()
			return true
		case 's':
			p.App().Volup()
			return true
		// speed
		case 'z':
			p.App().Spddown()
			return true
		case 'x':
			p.App().Spdup()
			return true
		// switch other panel
		case 'f':
			p.App().ShowLdpanel()
			return true
		}
	}
	return p.Panel.HandleEvent(ev)
}

// Draw is
func (p *PunosPanel) Draw() {
	p.update()
	p.Panel.Draw()
}

// mini logo
//		" _ __  _   _ _ __   ___  ___",
//		"| '_ \\| | | | '_ \\ / _ \\/ __|",
//		"| |_) | |_| | | | | (_) \\__ \\",
//		"| .__/ \\__,_|_| |_|\\___/|___/",
//		"|_|",

func (p *PunosPanel) update() {
	status, waveform := p.App().Status()
	base := []string{
		"	________   ___  ___   ________    ________   ________      ",
		"	|\\   __  \\ |\\  \\|\\  \\ |\\   ___  \\ |\\   __  \\ |\\   ____\\     ",
		"	\\ \\  \\|\\  \\\\ \\  \\\\\\  \\\\ \\  \\\\ \\  \\\\ \\  \\|\\  \\\\ \\  \\___|_    ",
		"	 \\ \\   ____\\\\ \\  \\\\\\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\\\ \\_____  \\   ",
		"	  \\ \\  \\___| \\ \\  \\\\\\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\\\|____|\\  \\  ",
		"	   \\ \\__\\     \\ \\_______\\\\ \\__\\\\ \\__\\\\ \\_______\\ ____\\_\\  \\ ",
		"	    \\|__|      \\|_______| \\|__| \\|__| \\|_______||\\_________\\",
		"	                                                \\|_________|",

		status["title"],
		status["position"],
		status["info"],
	}
	base = append(base, waveform...)

	p.text.SetLines(base)
}

// Init return just text box
func (p *PunosPanel) Init(app *App) {
	p.Panel.Init(app)

	p.SetTitle("Turn Table")
	p.text = views.NewTextArea()
	p.text.SetStyle(StyleNormal)
	p.text.SetLines([]string{
		//		" _ __  _   _ _ __   ___  ___",
		//		"| '_ \\| | | | '_ \\ / _ \\/ __|",
		//		"| |_) | |_| | | | | (_) \\__ \\",
		//		"| .__/ \\__,_|_| |_|\\___/|___/",
		//		"|_|",

		"	________   ___  ___   ________    ________   ________      ",
		"	|\\   __  \\ |\\  \\|\\  \\ |\\   ___  \\ |\\   __  \\ |\\   ____\\     ",
		"	\\ \\  \\|\\  \\\\ \\  \\\\\\  \\\\ \\  \\\\ \\  \\\\ \\  \\|\\  \\\\ \\  \\___|_    ",
		"	 \\ \\   ____\\\\ \\  \\\\\\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\\\ \\_____  \\   ",
		"	  \\ \\  \\___| \\ \\  \\\\\\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\\\|____|\\  \\  ",
		"	   \\ \\__\\     \\ \\_______\\\\ \\__\\\\ \\__\\\\ \\_______\\ ____\\_\\  \\ ",
		"	    \\|__|      \\|_______| \\|__| \\|__| \\|_______||\\_________\\",
		"	                                                \\|_________|",
	})
	p.SetContent(p.text)
	p.SetKeys([]string{"[ESC] Quit", "[SPACE] Play/Pause", "[Q/W] Rewind/Fastforward", "[A/S] Volume-/+", "[Z/X] Speed-/+", "[f] Switch Load Panel"})
}

//NewPunosPanel return PunosPanel
func NewPunosPanel(app *App) *PunosPanel {
	app.Logf("NewPunosPanel")
	p := &PunosPanel{}

	p.Init(app)
	return p
}
