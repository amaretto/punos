package ui

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/faiface/beep"
)

// Waveform has information of wave thined out
type Waveform struct {
	Wave           []int
	SampleInterval int
	WindowSize     int
	HeightMax      int
	ValMax         float64
	WaveDirPath    string
}

// GenRawWave generate raw waveform values([]float64) from volume values
func GenRawWave(streamer beep.StreamSeeker, sampleInterval int) []float64 {
	var tmp [2][2]float64
	var count, ncount int
	var rwave []float64
	// ToDo : currently, raw wave limit is 100000
	rwave = make([]float64, 100000)

	for {
		// check EOF
		if sn, sok := streamer.Stream(tmp[:1]); sn == 0 && !sok {
			break
		}
		samplel := tmp[0][0]
		sampler := tmp[0][1]

		sumSquare := math.Pow(samplel, 2)
		sumSquare += math.Pow(sampler, 2)
		value := math.Sqrt(sumSquare)

		if count%sampleInterval == 0 {
			rwave[ncount] = value
			ncount++
		}

		count++
	}

	rwave = rwave[:ncount]
	return rwave
}

// SmoothRawWave make raw wave values smoothly for visualizaiton
func SmoothRawWave(rwave []float64) {
	var sample float64
	sample = 3
	var sum float64
	for i := 0; i < len(rwave); i++ {
		if i < len(rwave)-int(sample) {
			sum = 0
			for j := 0; j < int(sample); j++ {
				sum += rwave[i+j]
			}
			rwave[i] = sum / sample
		} else {
			rwave[i] = rwave[i-1]
		}
	}
}

// NormalizeRawWave arrange wave values utilizing heightMax as height
func NormalizeRawWave(rwave []float64, heightMax, valMax float64) []int {
	var max float64
	var limit float64
	max = 1.0
	limit = heightMax

	var r []int
	r = make([]int, len(rwave))
	for i, num := range rwave {
		r[i] = int(math.Ceil(limit * num / max))
	}

	return r
}

// GetWave return part of wave by position
func GetWave(waves []int, pos, sampleInterval, windowSize int) []int {
	var wave []int
	wave = make([]int, windowSize)
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

// Wave2str generate strings express waveform
func Wave2str(wave []int, limit int) []string {
	var waveStr []string
	var fill []bool
	waveStr = make([]string, limit)
	fill = make([]bool, len(wave))

	for i := limit; i > 0; i-- {
		str := ""
		for j, num := range wave {
			if j == len(wave)/2-1 {
				str = str + "|"
			} else if num == i || fill[num] {
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

// GetWaveStr generate strings express waveform
func (wf *Waveform) GetWaveStr(pos int) []string {
	return Wave2str(GetWave(wf.Wave, pos, wf.SampleInterval, wf.WindowSize), wf.HeightMax)
}

// WriteWave write waveform to file
func WriteWave(wave []int, waveDirPath, title string) {
	// []int to []string
	var wsa []string // wave string array
	wsa = make([]string, len(wave))
	for i, num := range wave {
		wsa[i] = strconv.Itoa(num)
	}

	// out test
	wfile, err := os.Create(waveDirPath + "/" + title + ".txt")
	if err != nil {
		report(err)
	}
	defer wfile.Close()
	wfile.Write(([]byte)(strings.Join(wsa, " ")))
}

// LoadWave load waveform from file
func LoadWave(waveDirPath, title string) []int {
	// load test
	lfile, err := os.Open(waveDirPath + "/" + title + ".txt")
	if err != nil {
		report(err)
	}
	defer lfile.Close()

	var wave []int
	wave = make([]int, 1000000)
	var count int

	sc := bufio.NewScanner(lfile)
	// split by " "
	sc.Split(bufio.ScanWords)
	for sc.Scan() {
		wave[count], _ = strconv.Atoi(sc.Text())
		count++
	}

	return wave[:count]
}
