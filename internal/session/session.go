package session

import (
	"time"
	"fmt"
)

type Session struct {
	IsWork bool
	Duration time.Duration
	TotalPomodoros int
	IsCompleted bool
	StartTime time.Time
	EndTime time.Time
	SessionNumber int
	IsLongBreak bool
}

func (s *Session) StartWork() {
	if !s.IsWork || s.IsCompleted {
		s.SessionNumber++
	}
	s.IsWork = true
	s.IsCompleted = false
	s.StartTime = time.Now()
	s.EndTime = s.StartTime.Local().Add(s.Duration)
}

func (s *Session) StartBreak() {
	s.IsWork = false
	s.IsCompleted = false
	s.IsLongBreak = false
	s.StartTime = time.Now()
	s.EndTime = s.StartTime.Local().Add(s.Duration)
}

func (s *Session) EndSession() {
	if s.IsWork {
		s.TotalPomodoros++
	}
	s.IsWork = false
	s.IsCompleted = true
	s.EndTime = time.Now()
}

func (s *Session) ResetSession() {
	s.SessionNumber = 0
	s.TotalPomodoros = 0

}

func (s *Session) IsSessionActive() bool {
	return (time.Now().After(s.StartTime) && time.Now().Before(s.EndTime)) && !s.IsCompleted
}

func (s *Session) IsSessionCompleted() bool {
	return s.IsCompleted
}

func (s *Session) Summary() {
	fmt.Println("Session Summary:")
	fmt.Println("Session Duration:", s.Duration)
	fmt.Println("Session Number:", s.SessionNumber)
	fmt.Println("Total Pomodoros:", s.TotalPomodoros)
}