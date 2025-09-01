package audio

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// Plays a notification sound from a given file path
func PlayNotification(path string) error {
	fmt.Printf("Attempting to play audio file %s\n", path)
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %v\n", err)
		return err
	}
	defer file.Close()

	fmt.Println("File opened successfully")
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return err
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	volume := &effects.Volume{
		Streamer: streamer,
		Base: 2,
		Volume: -3.0,
		Silent: false,
	}
	
	done := make(chan bool)
	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		done <- true
	})))

	<- done
	return nil
}