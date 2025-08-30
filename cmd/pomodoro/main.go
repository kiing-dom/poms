package main
import (
	"fmt"
	"github.com/kiing-dom/poms/internal/timers"
	"github.com/kiing-dom/poms/internal/session"
)

func main() {
	s := &session.Session{}

	s.StartWork()
	fmt.Println("Starting Work Session:", s.SessionNumber)
	timers.Countdown(s.Duration, "Work")
	s.EndSession()
	fmt.Println("Work Session Complete. Good Job!")

	s.StartBreak()	
	fmt.Println("Starting Break Session:", s.SessionNumber)
	timers.Countdown(s.Duration, "Break")
	s.EndSession()
	fmt.Println("Break completed! Back to work")

	s.Summary()
}

