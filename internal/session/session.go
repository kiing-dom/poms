package session

import (
	"fmt"
	"time"
)

type Session struct {
	IsWork         bool
	WorkDuration   time.Duration
	BreakDuration  time.Duration
	TotalPomodoros int
	IsCompleted    bool
	StartTime      time.Time
	EndTime        time.Time
	SessionNumber  int
	IsLongBreak    bool
	IsPaused       bool
	PausedAt       time.Time
	TotalPaused    time.Duration
}

func (s *Session) StartWork() {
	if !s.IsWork || s.IsCompleted {
		s.SessionNumber++
	}
	s.IsWork = true
	s.IsCompleted = false
	s.StartTime = time.Now()
	s.EndTime = s.StartTime.Local().Add(s.WorkDuration)
	s.TotalPaused = 0
}

func (s *Session) StartBreak() {
	if s.IsWork {
		s.TotalPomodoros++
	}
	s.IsWork = false
	s.IsCompleted = false
	s.IsLongBreak = false
	s.StartTime = time.Now()
	s.EndTime = s.StartTime.Local().Add(s.BreakDuration)
	s.TotalPaused = 0
}

func (s *Session) EndSession() {
	s.IsWork = false
	s.IsCompleted = true
	s.EndTime = time.Now()
}

func (s *Session) ResetSession() {
	s.SessionNumber = 0
	s.TotalPomodoros = 0

}

func (s *Session) IsSessionActive() bool {
	return (time.Now().After(s.StartTime) && time.Now().Before(s.EndTime)) && !s.IsCompleted && !s.IsPaused
}

func (s *Session) IsSessionCompleted() bool {
	return s.IsCompleted
}

func (s *Session) Summary() {
	fmt.Println("Session Summary:")
	fmt.Println("Session Duration:", s.WorkDuration+s.BreakDuration)
	fmt.Println("Session Number:", s.SessionNumber)
	fmt.Println("Total Pomodoros:", s.TotalPomodoros)
}

func (s *Session) GetCurrentDuration() time.Duration {
	if s.IsWork {
		return s.WorkDuration
	}

	return s.BreakDuration
}

func (s *Session) Pause() {
	if !s.IsPaused && s.IsSessionActive() {
		s.IsPaused = true
		s.PausedAt = time.Now()
	}
}

func (s *Session) Resume() {
	if s.IsPaused {
		pauseDuration := time.Since(s.PausedAt)
		s.TotalPaused += pauseDuration
		s.EndTime = s.EndTime.Add(pauseDuration)
		s.IsPaused = false
	}
}
