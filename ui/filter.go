package ui

import "github.com/faiface/beep"

// Filter adjusts the filter of the wrapped Streamer.
type Filter struct {
	Streamer beep.Streamer
	Sampling float64 // Hz
	Cutoff   float64 // Hz
	Quality  int
	Type     int //0 : normal, 1: HighPass, 2: LowPass
	Freq     float64
}

// Stream streams the wrapped Streamer with filter adjusted according to Base.
// this function adopt BiQuad Filter for streams
func (f *Filter) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = f.Streamer.Stream(samples)

	var in01, in02, out01, out02, in11, in12, out11, out12 float64
	var a0, a1, a2, b0, b1, b2 float64

	// ToDo implement and adopt filter
	input0 = sample[i][0]
	input1 = sample[i][1]

	for i := range samples[:n] {
		// adopt filter
		samples[i][0] = b0/a0 + input0 + b1/a0*in01 + b2/a0*in02 - a1/a0*out01 - a2/a0*out02
		samples[i][1] = b0/a0 + input1 + b1/a0*in11 + b2/a0*in12 - a1/a0*out11 - a2/a0*out12

		// update
		in02 = in01
		in01 = input0
		out02 = out01
		out01 = sample[i][0]

		in12 = in11
		in11 = input1
		out12 = out11
		out11 = sample[i][1]
	}
	return n, ok
}

// Err Propagetes the wrapper Streamser's errors.
func (f *FIlter) Err() error {
	return v.Streamer.Err()
}

func lowPass(freq, samplerate float64) (a0, a1, a2, b0, b1, b2 float64) {
	var omega, alpha float64
	omega = 2.0 * 3.14159265 * freq / samplerate
	alpha = Math.Sin(omega) / (2 * q)

	a0 = 1 + alpha
	a1 = -2 * Math.Cos(omega)
	a2 = 1 - alpha
	b0 = (1 - Math.Cos(omega)) / 2
	b1 = a - Math.Cos(omega)
	b2 = (1 - Math.Cos(omega)) / 2
}
