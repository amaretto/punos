package ui

import (
	"math"

	"github.com/faiface/beep"
)

// GenWave is
func GenWave(streamer beep.StreamSeeker, sampleInterval int) []float64 {
	var tmp [2][2]float64
	var count, ncount int
	var wave []float64
	wave = make([]float64, 100000)

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

		if count%800 == 0 {
			wave[ncount] = value
			ncount++
		}

		count++
	}

	wave = wave[:ncount]
	return wave
}

// Smooth is
func Smooth(nums []float64) {
	var sample float64
	sample = 3
	var sum float64
	for i := 0; i < len(nums); i++ {
		if i < len(nums)-int(sample) {
			sum = 0
			for j := 0; j < int(sample); j++ {
				sum += nums[i+j]
			}
			nums[i] = sum / sample
		} else {
			nums[i] = nums[i-1]
		}
	}
}

// Normalize is
func Normalize(nums []float64, heightMax, valMax float64) []int {
	var max float64
	var limit float64
	max = 1.0
	limit = heightMax

	var r []int
	r = make([]int, len(nums))
	for i, num := range nums {
		r[i] = int(math.Ceil(limit * num / max))
	}

	return r
}

// GetWave is
func GetWave(waves []int, pos, sampleRate, windowSize int) []int {
	var wave []int
	wave = make([]int, windowSize)
	center := int(pos / 800)

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

// Wave2str is
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