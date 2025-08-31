package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/kiing-dom/poms/internal/session"
)

type Model struct {
	session *session.Session
	timer *time.Timer
	running bool
}

func NewModel(session *session.Session) Model {
	return Model {
		session: session,
		timer: nil,
		running: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w":
			m.session.StartWork()
			m.running = true
			return m, m.startTimer()
		case "b":
			m.session.StartBreak()
			m.running = true
			return m, m.startTimer()
		case "e":
			m.session.EndSession()
			m.running = false
			if m.timer != nil {
				m.timer.Stop()
			}
			return m, nil
		case "q":
			return m, tea.Quit
		}
	case tickMsg:
		if m.running && m.session.IsSessionActive() {
			return m, m.startTimer()
		}
		m.running = false
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	status := "Idle"
	if m.running {
		if m.session.IsWork {
			status = "Working"
		} else {
			status = "On Break"
		}
	}

	return fmt.Sprintf("Poms\nStatus: %s\nPomodoros: %d\nSession: %d\n", status, m.session.TotalPomodoros, m.session.SessionNumber)
}

type tickMsg time.Time

func (m Model) startTimer() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}