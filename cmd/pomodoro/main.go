package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kiing-dom/poms/internal/cli"
	"github.com/kiing-dom/poms/internal/session"
	"github.com/kiing-dom/poms/internal/timers"
	tui "github.com/kiing-dom/poms/internal/tui"
)

func main() {
	config := cli.ParseFlags()

	if config.ShouldUseTUI() {
		workDuration := time.Duration(config.WorkMinutes) * time.Minute
		breakDuration := time.Duration(config.BreakMinutes) * time.Minute
		if err := StartTUI(workDuration, breakDuration); err != nil {
			fmt.Printf("Error starting TUI: %v\n", err)
		}
	} else {
		runCLIMode(config)
	}
}

func StartTUI(workDuration, breakDuration time.Duration) error {
	s := &session.Session{Duration: workDuration}
	m := tui.NewModel(s)
	p := tea.NewProgram(&m)
	_, err := p.Run()
	return err
}

func runCLIMode(config cli.Config) {
	s := &session.Session{}

	s.Duration = time.Duration(config.WorkMinutes) * time.Minute
	s.StartWork()
	fmt.Println("Starting Work Session:", s.SessionNumber)
	timers.Countdown(s.Duration, "Work")
	s.EndSession()
	fmt.Println("Work Session Complete. Good Job!")

	s.Duration = time.Duration(config.BreakMinutes) * time.Minute
	s.StartBreak()
	fmt.Println("Starting Break Session:", s.SessionNumber)
	timers.Countdown(s.Duration, "Break")
	s.EndSession()
	fmt.Println("Break over. Back to Work!")

	s.Summary()
}
