package timers

import (
	"fmt"
	"time"
)

func Countdown(duration time.Duration, label string) {
	durationSeconds := int(duration.Seconds())

	for i := durationSeconds; i > 0; i-- {
		fmt.Println(label, ":", i, "seconds left...")
		time.Sleep(1 * time.Second)
	}
}