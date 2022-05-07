package player

import (
	"strings"
)

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

func (w *WaveformPanel) update(wave []byte, pos int) {
	_, _, width, height := w.GetRect()
	w.SetText(strings.Join(getWaveStr(wave, pos, width-2, height-2), "\n"))
}

// GetWave return part of wave by position
func getWave(waves []byte, pos, sampleInterval, windowSize int) []byte {
	var wave []byte
	wave = make([]byte, windowSize)
	center := int(pos / sampleInterval)

	var idx int
	for i := 0; i < windowSize; i++ {
		idx = center - (windowSize/2 - 1) + i
		// front
		if idx < 0 {
			continue
		}
		// back
		if idx >= len(waves) {
			break
		}
		wave[i] = waves[idx]
	}
	return wave
}

// wave2str generate strings express waveform
func wave2str(wave []byte, limit int) []string {
	var waveStr []string
	var fill []bool
	waveStr = make([]string, limit)
	fill = make([]bool, len(wave))

	for i := limit; i > 0; i-- {
		str := ""
		for j, num := range wave {
			if j == len(wave)/2-1 {
				str = str + "|"
			} else if int(num) >= i || fill[num] {
				str = str + "#"
				fill[num] = true
			} else {
				str = str + " "
			}
		}
		waveStr[limit-i] = str
	}
	return waveStr
}

// getWaveStr generate strings express waveform
func getWaveStr(wave []byte, pos, width, height int) []string {
	// magic number
	sampleInterval := 800
	return wave2str(getWave(wave, pos, sampleInterval, width), height)
}
