package player

type ProgressBar struct {
	*DefaultView
}

func NewProgressBar() *ProgressBar {
	p := &ProgressBar{
		DefaultView: NewDefaultView("Progress"),
	}
	p.initProgressBar()
	return p
}

func (p *ProgressBar) initProgressBar() {
	p.SetText("no music selected")
}
