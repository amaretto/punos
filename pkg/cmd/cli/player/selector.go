package player

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Selector is panel for selecting music
type Selector struct {
	*tview.Flex
	helpModal    tview.Primitive
	analyzeModal tview.Primitive

	player *Player

	musicListView *tview.Table
}

func newSelector(player *Player) *Selector {
	s := &Selector{
		player: player,
		Flex:   tview.NewFlex(),

		musicListView: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
	}
	s.SetTitle("selector")

	// set header
	headers := []string{
		"Status",
		"Album",
		"Authors",
		"Title",
		"Duration",
		"BPM",
	}
	for i, header := range headers {
		s.musicListView.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold,
		})
	}

	keyBindingHelp := [][]string{
		{"Key", "Description"},
		{"Esc", "Quit"},
		{"d", "SelectMusic"},
		{"j", "Down Cursor"},
		{"k", "Up Cursor"},
		{"a", "Analyze All Music"},
		{"l", "Move Selected Dir"},
		{"h", "Move Parent Dir"},
	}
	s.helpModal = newHelpModal(keyBindingHelp)
	s.analyzeModal = newMsgModal(30, 5, "hogehogehogehoge")

	s.musicListView.SetBorder(true).SetTitleAlign(tview.AlignLeft).SetTitle("MusicList")
	s.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(s.musicListView, 0, 4, false),
			0, 1, false)
	s.SetKeyHandler()

	s.update()
	return s
}

func (s *Selector) update() {
	s.musicListView.Clear()
	// music list udpate
	for i, musicInfo := range s.player.musics.List {
		s.musicListView.SetCell(i+1, 0, tview.NewTableCell(musicInfo.Status).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 1, tview.NewTableCell(musicInfo.Album).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 2, tview.NewTableCell(musicInfo.Authors).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 3, tview.NewTableCell(musicInfo.Title).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 4, tview.NewTableCell(musicInfo.Duration).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 5, tview.NewTableCell(musicInfo.BPM).SetMaxWidth(1).SetExpansion(1))
	}

	// dir list update
	for i, dirPath := range s.player.musics.ListDirs() {
		s.musicListView.SetCell(i+1+len(s.player.musics.List), 0, tview.NewTableCell("Dir").SetMaxWidth(1).SetExpansion(1).SetTextColor(tcell.ColorBlue))
		s.musicListView.SetCell(i+1+len(s.player.musics.List), 3, tview.NewTableCell(dirPath).SetMaxWidth(1).SetExpansion(1).SetTextColor(tcell.ColorBlue))
	}
}

// SetKeyHandler is
func (s *Selector) SetKeyHandler() {
	s.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Rune() {
		case 'a':
			for _, m := range s.player.musics.List {
				if m.Status == "Not Analyzed" {
					s.player.analyzer.AnalyzeMusic(m)
				}
				// ToDo:this doesn't work........
				s.update()
			}
			//s.player.musics.ListMusics()
			s.update()
		case 'd':
			// select music
			row, _ := s.musicListView.GetSelection()
			if row <= len(s.player.musics.List) {
				s.player.LoadMusic(s.player.musics.List[row-1])
			}
		case 'h':
			// move parent dir
			s.player.musics.MoveParentDir()
			s.update()
			s.musicListView.Select(1, 1)
		case 'l':
			row, _ := s.musicListView.GetSelection()
			if row > len(s.player.musics.List) {
				s.player.musics.MoveChildDir(s.musicListView.GetCell(row, 3).Text)
			}
			s.update()
			s.musicListView.Select(1, 1)
		}
		return e
	})
}
