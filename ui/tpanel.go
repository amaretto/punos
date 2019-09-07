package ui

import (
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var (
	// StyleNormal is
	StyleNormal = tcell.StyleDefault.
			Foreground(tcell.ColorSilver).
			Background(tcell.ColorBlack)
	// StyleB2B is
	StyleB2B = tcell.StyleDefault.
			Foreground(tcell.ColorSilver).
			Background(tcell.ColorTeal)
	// StyleSync is
	StyleSync = tcell.StyleDefault.
			Foreground(tcell.ColorSilver).
			Background(tcell.ColorMaroon)
)

type trntblModel struct {
	m *TrntblPanel
}

// TrntblPanel is
type TrntblPanel struct {
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
func (t *TrntblPanel) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEsc:
			t.App().Quit()
			return true
		}

		switch unicode.ToLower(ev.Rune()) {
		case ' ':
			t.App().PlayPause()
			return true
		case 'c':
			t.App().Cue()
			return true
		case 'w':
			t.App().Fforward()
			return true
		case 'q':
			t.App().Rewind()
			return true
		// volume
		case 'a':
			t.App().Voldown()
			return true
		case 's':
			t.App().Volup()
			return true
		// speed
		case 'z':
			t.App().Spddown()
			return true
		case 'x':
			t.App().Spdup()
			return true
		// switch other panel
		case 'f':
			t.App().ShowLdpanel()
			return true
		// switch mode
		case 'n':
			t.App().SetMode("normal")
			return true
		case 'b':
			t.App().SetMode("b2b")
			return true
		case 'm':
			t.App().SetMode("sync")
			return true
		}
	}
	return t.Panel.HandleEvent(ev)
}

// Draw is
func (t *TrntblPanel) Draw() {
	t.update()
	t.Panel.Draw()
}

func (t *TrntblPanel) update() {

	// set Style
	if t.App().Mode == "normal" {
		t.text.SetStyle(StyleNormal)
	} else if t.App().Mode == "b2b" {
		t.text.SetStyle(StyleB2B)
	} else if t.App().Mode == "sync" {
		t.text.SetStyle(StyleSync)
	}

	status, waveform := t.App().Status()
	base := []string{
		" _ __  _   _ _ __   ___  ___",
		"| '_ \\| | | | '_ \\ / _ \\/ __|",
		"| |_) | |_| | | | | (_) \\__ \\",
		"| .__/ \\__,_|_| |_|\\___/|___/",
		"|_|",
		status["title"],
		status["position"],
		status["info"],
		//		status["volume"],
		//		status["speed"],
	}
	base = append(base, waveform...)

	t.text.SetLines(base)
}

// Init return just text box
func (t *TrntblPanel) Init(app *App) {
	t.Panel.Init(app)

	t.SetTitle("Turn Table")
	t.text = views.NewTextArea()
	t.text.SetStyle(StyleNormal)
	t.text.SetLines([]string{
		" _ __  _   _ _ __   ___  ___",
		"| '_ \\| | | | '_ \\ / _ \\/ __|",
		"| |_) | |_| | | | | (_) \\__ \\",
		"| .__/ \\__,_|_| |_|\\___/|___/",
		"|_|",
	})
	t.SetContent(t.text)
	t.SetKeys([]string{"[ESC] Quit", "[SPACE] Play/Pause", "[Q/W] Rewind/Fastforward", "[A/S] Volume-/+", "[Z/X] Speed-/+", "[f] Switch Load Panel"})
}

//NewTrntblPanel return TrntblPanel
func NewTrntblPanel(app *App) *TrntblPanel {
	app.Logf("NewTrntblPanel")
	t := &TrntblPanel{}

	//t.Panel.Init(app)
	//t.content = views.NewCellView()
	//t.SetContent(t.content)

	//t.content.SetModel(&trntblModel{t})
	//t.content.SetStyle(StyleNormal)

	//t.SetTitle("hoge")

	t.Init(app)
	return t
}
