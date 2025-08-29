package main
import (
	"fmt"
	"time"
	"github.com/kiing-dom/poms/internal/timers"
)

func main() {
	workDuration := 10 * time.Second
	breakDuration := 1 * time.Second

	timers.Countdown(workDuration, "Work")
	fmt.Println("Work block complete. Time to take a break")
	time.Sleep(2 * time.Second)
	timers.Countdown(breakDuration, "Break")
}

