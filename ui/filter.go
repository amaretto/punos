package ui

import (
	"fmt"
	"math"

	"github.com/faiface/beep"
)

// Filter adjusts the filter of the wrapped Streamer.
type Filter struct {
	Streamer beep.Streamer
	Sampling float64 // sample freq (Hz)
	Freq     float64 // cutoff freq (Hz)
	Q        float64
	Type     int //0 : normal, 1: HighPass, 2: LowPass
}

// Stream streams the wrapped Streamer with filter adjusted according to Base.
// this function adopt BiQuad Filter(HighPass/LowPass) for streams
func (f *Filter) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = f.Streamer.Stream(samples)

	var input0, input1 float64
	var in01, in02, out01, out02, in11, in12, out11, out12 float64
	var a0, a1, a2, b0, b1, b2 float64

	// ToDo : implemet other filter
	a0, a1, a2, b0, b1, b2 = lowPassParam(f.Sampling, f.Freq, f.Q)

	for i := range samples[:n] {
		input0 = samples[i][0]
		input1 = samples[i][1]

		// adopt filter
		samples[i][0] = (b0/a0)*input0 + (b1/a0)*in01 + (b2/a0)*in02 - (a1/a0)*out01 - (a2/a0)*out02
		samples[i][1] = (b0/a0)*input1 + (b1/a0)*in11 + (b2/a0)*in12 - (a1/a0)*out11 - (a2/a0)*out12

		fmt.Println("before:", input0, ",after:", samples[i][0])

		// update
		in02 = in01
		in01 = input0
		out02 = out01
		out01 = samples[i][0]

		in12 = in11
		in11 = input1
		out12 = out11
		out11 = samples[i][1]
	}
	return n, ok
}

// Err Propagetes the wrapper Streamser's errors.
func (f *Filter) Err() error {
	return f.Streamer.Err()
}

func lowPassParam(samplerate, freq, q float64) (a0, a1, a2, b0, b1, b2 float64) {
	var omega, alpha float64
	omega = 2.0 * 3.14159265 * freq / samplerate
	alpha = math.Sin(omega) / (2.0 * q)

	a0 = 1.0 + alpha
	a1 = -2.0 * math.Cos(omega)
	a2 = 1.0 - alpha
	b0 = (1.0 - math.Cos(omega)) / 2.0
	b1 = 1.0 - math.Cos(omega)
	b2 = (1.0 - math.Cos(omega)) / 2.0

	return
}

func highPassParam(samplerate, freq, q float64) (a0, a1, a2, b0, b1, b2 float64) {
	var omega, alpha float64
	omega = 2.0 * 3.14159265 * freq / samplerate
	alpha = math.Sin(omega) / (2.0 * q)

	a0 = 1.0 + alpha
	a1 = -2.0 * math.Cos(omega)
	a2 = 1.0 - alpha
	b0 = (1.0 - math.Cos(omega)) / 2.0
	b1 = -1 * (1.0 + math.Cos(omega))
	b2 = (1.0 + math.Cos(omega)) / 2.0

	return
}
