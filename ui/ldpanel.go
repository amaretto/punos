package ui

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

//var (
//	// StyleNormal is
//	StyleNormal = tcell.StyleDefault.
//		Foreground(tcell.ColorSilver).
//		Background(tcell.ColorBlack)
//)

type loadModel struct {
	l *LoadPanel
}

// LoadPanel is
type LoadPanel struct {
	width     int
	height    int
	curx      int
	cury      int
	content   *views.CellView
	text      *views.TextArea
	selected  string
	items     []string
	lines     []string
	styles    []tcell.Style
	musicDir  string
	musicList []string

	Panel
}

// HandleEvent is
func (l *LoadPanel) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEsc:
			l.App().Quit()
			return true
		}

		switch unicode.ToLower(ev.Rune()) {
		// Cursor
		case 'j':
			l.MoveCursor(0, 1)
			return true
		case 'k':
			l.MoveCursor(0, -1)
			return true
		// Load music
		case 'l':
			l.App().LoadMusic(l.selected)
			l.App().ShowPunosPanel()
			return true
		// swtich other panel
		case 'f':
			l.App().ShowPunosPanel()
			return true
			// analyze all music data
		case 'a':
			l.App().Analyze()
			return true
		}
	}
	return l.Panel.HandleEvent(ev)
}

//GetCell is method for Model
func (model *loadModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	var ch rune
	var style tcell.Style

	l := model.l

	if y < 0 || y >= len(l.lines) {
		return ch, StyleNormal, nil, 1
	}

	if x >= 0 && x < len(l.lines[y]) {
		ch = rune(l.lines[y][x])
	} else {
		ch = ' '
	}
	style = l.styles[y]
	if l.items[y] == l.selected {
		style = style.Reverse(true)
	}
	return ch, style, nil, 1
}

// GetBounds search boud of the application window
func (model *loadModel) GetBounds() (int, int) {
	l := model.l
	y := len(l.lines)
	x := 0
	for _, l := range l.lines {
		if x < len(l) {
			x = len(l)
		}
	}
	return x, y
}

// GetCursor return position of the cursor
func (model *loadModel) GetCursor() (int, int, bool, bool) {
	l := model.l
	return l.curx, l.cury, true, false
}

// MoveCursor move cursor
func (model *loadModel) MoveCursor(offx, offy int) {
	l := model.l
	l.curx += offx
	l.cury += offy
	l.updateCursor(true)
}

// SetCursor set cursor
func (model *loadModel) SetCursor(x, y int) {
	l := model.l
	l.curx = x
	l.cury = y
	l.updateCursor(true)
}

// MoveCursor move cursor
func (l *LoadPanel) MoveCursor(offx, offy int) {
	l.curx += offx
	l.cury += offy
	l.updateCursor(true)
}

func (l *LoadPanel) unselect() {
	l.cury = 0
	l.curx = 0
	l.updateCursor(false)
}

func (l *LoadPanel) updateCursor(selected bool) {
	// for bound
	if l.curx > l.width-1 {
		l.curx = l.width - 1
	}
	if l.cury > l.height-1 {
		l.cury = l.height - 1
	}
	if l.curx < 0 {
		l.curx = 0
	}
	if l.cury < 0 {
		l.cury = 0
	}
	if selected && l.height > 0 {
		if l.selected == "" {
			l.curx = 0
			l.cury = 0
		}
		l.selected = l.items[l.cury]
	} else {
		l.selected = ""
	}
}

// Draw is
func (l *LoadPanel) Draw() {
	l.update()
	l.Panel.Draw()
}

func (l *LoadPanel) update() {

	l.items = l.App().ListMusic()

	if sel := l.selected; sel != "" {
		l.selected = ""
		cury := 0
		for _, item := range l.items {
			if item == sel {
				l.selected = item
				l.cury = cury
			}
			cury++
		}
	}
	//	if err != nil {
	//		report(err)
	//	}

	lines := make([]string, 0, len(l.items))
	styles := make([]tcell.Style, 0, len(l.items))

	l.height = 0
	l.width = 0

	// ToDo : Refactoring
	for _, item := range l.items {
		//ToDo : validation item?
		line := fmt.Sprintf(item)
		if len(line) > l.width {
			l.width = len(line)
		}
		l.height++

		lines = append(lines, line)
		var style tcell.Style
		style = StyleNormal
		styles = append(styles, style)
	}

	l.lines = lines
	l.styles = styles

	//	l.text.SetLines([]string{
	//		" _ __  _   _ _ __   ___  ___",
	//		"| '_ \\| | | | '_ \\ / _ \\/ __|",
	//		"| |_) | |_| | | | | (_) \\__ \\",
	//		"| .__/ \\__,_|_| |_|\\___/|___/",
	//		"|_|",
	//		time.Now().String(),
	//	})
}

//NewLoadPanel return LoadPanel
func NewLoadPanel(app *App) *LoadPanel {
	app.Logf("NewLoadPanel")
	l := &LoadPanel{}
	l.Panel.Init(app)
	l.SetTitle("Load Music")

	//l.text = views.NewTextArea()
	//l.text.SetStyle(StyleNormal)
	//l.text.SetLines([]string{
	//	" _ __  _   _ _ __   ___  ___",
	//	"| '_ \\| | | | '_ \\ / _ \\/ __|",
	//	"| |_) | |_| | | | | (_) \\__ \\",
	//	"| .__/ \\__,_|_| |_|\\___/|___/",
	//	"|_|",
	//})
	//l.SetContent(l.text)

	l.content = views.NewCellView()
	l.content.SetModel(&loadModel{l})
	l.content.SetStyle(StyleNormal)
	l.SetContent(l.content)

	l.SetKeys([]string{"[ESC] Quit", "[J] Down", "[K] Up", "[L] load", "[A] Analyze All Music"})
	return l
}
