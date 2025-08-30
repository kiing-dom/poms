package main
import (
	"fmt"
	"time"
	"github.com/kiing-dom/poms/internal/timers"
	"github.com/kiing-dom/poms/internal/session"
	"github.com/kiing-dom/poms/internal/cli"
)

func main() {
	workMin, breakMin := cli.ParseFlags()
	s := &session.Session{}
	
	s.Duration = time.Duration(workMin) * time.Minute
	s.StartWork()
	fmt.Println("Starting Work Session:", s.SessionNumber)
	timers.Countdown(s.Duration, "Work")
	s.EndSession()
	fmt.Println("Work Session Complete. Good Job!")

	s.Duration = time.Duration(breakMin) * time.Minute
	s.StartBreak()	
	fmt.Println("Starting Break Session:", s.SessionNumber)
	timers.Countdown(s.Duration, "Break")
	s.EndSession()
	fmt.Println("Break completed! Back to work")

	s.Summary()
}

