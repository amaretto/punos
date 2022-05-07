package player

import (
	"os"
	"regexp"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Selector is panel for selecting music
type Selector struct {
	*tview.Flex
	player *Player

	musicListView *tview.Table
	musicDetail   *DefaultView
}

func newSelector(player *Player) *Selector {
	s := &Selector{
		player: player,
		Flex:   tview.NewFlex(),

		musicListView: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		musicDetail:   NewDefaultView("Music Detail"),
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
	for i, musicInfo := range s.player.musics.List {
		s.musicListView.SetCell(i+1, 0, tview.NewTableCell(musicInfo.Status).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 1, tview.NewTableCell(musicInfo.Album).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 2, tview.NewTableCell(musicInfo.Authors).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 3, tview.NewTableCell(musicInfo.Title).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 4, tview.NewTableCell(musicInfo.Duration).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 5, tview.NewTableCell(musicInfo.BPM).SetMaxWidth(1).SetExpansion(1))
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

// SetKeyHandler is
func (s *Selector) SetKeyHandler() {
	s.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Rune() {
		case 'a':
			for _, m := range s.player.musics.List {
				if m.Status == "Not Analyzed" {
					s.player.analyzer.AnalyzeMusic(m)
				}
			}
		case 'l':
			// load music
			row, _ := s.musicListView.GetSelection()
			s.player.LoadMusic(s.player.musics.List[row-1])
		}
		return e
	})
}

func (s *Selector) listMusic(musicPath string) []string {
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
