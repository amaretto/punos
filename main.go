package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	fmt.Println("\n\n")
	fmt.Println("_ __  _   _ _ __   ___  ___")
	fmt.Println("| '_ \\| | | | '_ \\ / _ \\/ __|")
	fmt.Println("| |_) | |_| | | | | (_) \\__ \\")
	fmt.Println("| .__/ \\__,_|_| |_|\\___/|___/")
	fmt.Println("|_|")
	fmt.Println("\n\n")

	var in string
	fmt.Printf("[empty]>>")
	fmt.Scanln(&in)

	f, err := os.Open("mp3/01.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done

	fmt.Println(in)
}
