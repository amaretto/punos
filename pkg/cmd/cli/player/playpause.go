package player

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const pauseString = `



 _ __   __ _ _   _ ___  ___ 
| '_ \ / _' | | | / __|/ _ \
| |_) | (_| | |_| \__ \  __/
| .__/ \__,_|\__,_|___/\___|
|_|                         
`

const playingString = `



                             _             _                    
 _ __   _____      __  _ __ | | __ _ _   _(_)_ __   __ _        
| '_ \ / _ \ \ /\ / / | '_ \| |/ _' | | | | | '_ \ / _' |       
| | | | (_) \ V  V /  | |_) | | (_| | |_| | | | | | (_| |       
|_| |_|\___/ \_/\_/   | .__/|_|\__,_|\__, |_|_| |_|\__, |       
                      |_|            |___/         |___/        
`

type PlayPausePanel struct {
	*DefaultView
}

func NewPlayPausePanel() *PlayPausePanel {
	p := &PlayPausePanel{
		DefaultView: NewDefaultView("PlayPause"),
	}
	p.initPlayPausePanel()
	p.SetTextAlign(tview.AlignCenter)
	return p
}

func (p *PlayPausePanel) setPause() {
	p.SetText(pauseString)
	p.SetTextColor(tcell.ColorYellow)

}

func (p *PlayPausePanel) setPlaying() {
	_, _, width, _ := p.GetRect()
	p.SetText(shiftAndFitLogo(playingString, width-2))
	p.SetTextColor(tcell.ColorGreen)
}

//ToDo: Fit height
func shiftAndFitLogo(logo string, width int) string {
	logoStrings := strings.Split(logo, "\n")
	var maxLen int
	for _, s := range logoStrings {
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}
	shiftNum := (time.Now().Unix() * 10) % int64(maxLen)
	for i := 0; i < len(logoStrings); i++ {
		logoString := []byte(logoStrings[i])
		if len(logoString) == 0 {
			continue
		}
		// shift every 1sec
		logoString = append(logoString[shiftNum:], logoString[:shiftNum]...)
		// fit
		if len(logoString) > width {
			logoString = logoString[:width]
		}
		logoStrings[i] = string(logoString)
	}
	return strings.Join(logoStrings, "\n")
}

func (p *PlayPausePanel) update(isPlay bool) {
	if isPlay {
		p.setPlaying()
		return
	}
	p.setPause()
	return
}

func (p *PlayPausePanel) initPlayPausePanel() {
	p.setPause()
}
