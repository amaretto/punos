package player

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"

	mdl "github.com/amaretto/punos/pkg/cmd/cli/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

// Selector is panel for selecting music
type Selector struct {
	*tview.Flex
	player *Player

	musicListView *tview.Table
	musicDetail   *DefaultView
	musicList     []*mdl.MusicInfo
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
	// list music file path from path
	musicPathList := s.listMusic("dummy/path")
	s.musicList = make([]*mdl.MusicInfo, 0)

	// get music info from db
	dbPath := "mp3/test.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		report(err)
	}
	rows, err := db.Query("SELECT path, title, album, authors, duration, sampleRate, format, bpm FROM music")
	if err != nil {
		report(err)
	}

	for rows.Next() {
		mi := &mdl.MusicInfo{}

		if err := rows.Scan(&mi.Path, &mi.Title, &mi.Album, &mi.Authors, &mi.Duration, &mi.SampleRate, &mi.Format, &mi.BPM); err != nil {
			report(err)
		}

		if contains(mi.Path, musicPathList) {
			mi.Status = "✔"
			musicPathList = del(mi.Path, musicPathList)
		} else {
			mi.Status = "Moved"
		}

		s.musicList = append(s.musicList, mi)
	}

	for _, musicPath := range musicPathList {
		mi := &mdl.MusicInfo{}
		mi.Path = musicPath
		mi.Title = filepath.Base(musicPath)
		mi.Status = "Not Analyzed"
		s.musicList = append(s.musicList, mi)
	}

	for i, musicInfo := range s.musicList {
		s.musicListView.SetCell(i+1, 0, tview.NewTableCell(musicInfo.Status).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 1, tview.NewTableCell(musicInfo.Album).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 2, tview.NewTableCell(musicInfo.Authors).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 3, tview.NewTableCell(musicInfo.Title).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 4, tview.NewTableCell(musicInfo.Duration).SetMaxWidth(1).SetExpansion(1))
		s.musicListView.SetCell(i+1, 5, tview.NewTableCell(musicInfo.BPM).SetMaxWidth(1).SetExpansion(1))
	}

	s.musicListView.SetBorder(true).SetTitleAlign(tview.AlignLeft).SetTitle("MusicList")

	s.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(s.musicListView, 0, 4, false),
			//.AddItem(s.musicDetail, 0, 2, false),
			0, 1, false)

	s.SetKeyHandler()
	return s
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
			logrus.Debug(s.musicList)
			for _, m := range s.musicList {
				if m.Status == "Not Analyzed" {
					s.player.analyzer.AnalyzeMusic(m)
				}
			}
		case 'l':
			// load music
			row, _ := s.musicListView.GetSelection()
			s.player.LoadMusic(s.musicList[row-1])
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
