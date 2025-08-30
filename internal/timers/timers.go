package timers

import (
	"fmt"
	"time"
)

func Countdown(duration time.Duration, label string) {
	durationSeconds := int(duration.Seconds())

	for i := durationSeconds; i > 0; i-- {
		min := i / 60
		sec := i % 60
		fmt.Printf("\r%s: %02d:%02d", label, min, sec)
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
}