package player

const logo = `
________   ___  ___   ________    ________   ________      
|\   __ \ |\  \|\  \ |\   ___  \ |\   __  \ |\   ____\     
\ \  \|\ \ \ \  \\  \\ \  \\ \  \\ \  \|\  \\ \  \___|_    
 \ \   ___\ \ \  \\  \\ \  \\ \  \\ \  \\\  \\ \_____  \   
  \ \  \___| \ \  \\  \\ \  \\ \  \\ \  \\\  \\|____|\  \  
   \ \__\     \ \_______\\\__\\ \__\\ \_______\ ____\_\  \ 
    \|__|      \|_______| \|__| \|__| \|_______||\________\
                                                \|________|
`

type WaveformPanel struct {
	*DefaultView
}

func NewWaveformPanel() *WaveformPanel {
	w := &WaveformPanel{
		DefaultView: NewDefaultView("Waveform"),
	}
	w.initWaveformPanel()
	return w
}

func (w *WaveformPanel) initWaveformPanel() {
	w.SetText(logo)
}
