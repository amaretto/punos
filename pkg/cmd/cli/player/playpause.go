package player

type PlayPausePanel struct {
	*DefaultView
}

func NewPlayPausePanel() *PlayPausePanel {
	p := &PlayPausePanel{
		DefaultView: NewDefaultView("PlayPause"),
	}
	p.initPlayPausePanel()
	return p
}

func (p *PlayPausePanel) initPlayPausePanel() {
	p.SetText("ToDo:show animation")
}
