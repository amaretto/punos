package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/faiface/beep/mp3"
	//gomp3 "github.com/hajimehoshi/go-mp3"
)

func main() {

	f, err := os.Open("mp3/02.mp3")
	if err != nil {
		report(err)
	}
	streamer, _, err := mp3.Decode(f)

	cf, err := os.OpenFile("hoge.csv", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer cf.Close()
	err = cf.Truncate(0)
	if err != nil {
		log.Fatal("Error:", err)
	}

	writer := csv.NewWriter(cf)

	var tmp [2][2]float64
	var count int
	for {
		// EOF
		if sn, sok := streamer.Stream(tmp[:1]); sn == 0 && !sok {
			break
		}
		samplel := tmp[0][0]
		sampler := tmp[0][1]

		sumSquare := math.Pow(samplel, 2)
		sumSquare += math.Pow(sampler, 2)
		value := math.Sqrt(sumSquare)

		posstr := fmt.Sprint(streamer.Position())
		valstr := fmt.Sprint(value)

		writer.Write([]string{posstr, valstr})

		count++
		//fmt.Printf("pos:%v,l:%v, r:%v\n", streamer.Position(), samplel, sampler)

		//fmt.Printf("pos:%v,value:%v\n", streamer.Position(), value)
		if count == 1000000 {
			break
		}
	}
	writer.Flush()
	//fmt.Println(count)
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
