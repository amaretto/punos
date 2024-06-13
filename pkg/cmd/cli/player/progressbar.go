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

func (p *ProgressBar) update(position, length int) {
	_, _, width, _ := p.GetRect()
	pb := genProgressBar(width-4, position, length)
	p.SetText(pb)
}

func genProgressBar(width, position, length int) string {
	current := float64(width) * float64(position) / float64(length)
	str := "["
	for i := 0; i < width; i++ {
		if i > int(current) {
			str = str + "-"
		} else {
			str = str + "▇"
		}
	}
	str = str + "]"
	return str
}
