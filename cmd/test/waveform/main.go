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
	var nums []float64
	nums = make([]float64, 100000)
	ncount := 0

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
		//
		//		posstr := fmt.Sprint(streamer.Position())
		//		valstr := fmt.Sprint(value)

		if count%800 == 0 {
			//if count >= 900000 && count%200 == 0 {
			nums[ncount] = value
			ncount++
			//			writer.Write([]string{posstr, valstr})
		}

		count++
		//fmt.Printf("pos:%v,l:%v, r:%v\n", streamer.Position(), samplel, sampler)

		//fmt.Printf("pos:%v,value:%v\n", streamer.Position(), value)
		//		if count == 1000000 {
		//			break
		//		}
	}
	//	writer.Flush()

	nums = nums[:ncount]

	smooth(nums)
	smooth(nums)
	smooth(nums)
	smooth(nums)

	nnums := normalize(nums)
	nnums = nnums[3000:5000]

	for i, num := range nnums {
		istr := fmt.Sprint(i)
		numstr := fmt.Sprint(num)
		writer.Write([]string{istr, numstr})
	}
	writer.Flush()

	printwave(nnums)
	//fmt.Println(count)
}

func smooth(nums []float64) {
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

func normalize(nums []float64) []int {
	var max float64
	var limit float64
	max = 1.0
	limit = 30.0

	var r []int
	r = make([]int, len(nums))
	for i, num := range nums {
		r[i] = int(math.Ceil(limit * num / max))
	}

	return r
}

func printwave(nums []int) {
	//var limit int
	//limit = 15
	//	var out []string
	//	out = make([]string, limit)
	file, err := os.Create("fuga.txt")
	if err != nil {
		report(err)
	}
	defer file.Close()

	//	var gage []byte
	//	gage = make([]byte, limit)
	//	for i := 0; i < limit; i++ {
	//		gage[i] = '#'
	//	}
	//
	gage := "##############################\n"

	for _, num := range nums {
		if num != 0 {
			file.WriteString(gage[31-num:])
			//file.Write("\n")
		}
	}
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
