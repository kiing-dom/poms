package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kiing-dom/poms/internal/session"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B6B")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B"))

	statusStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Margin(1, 0)

	workStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true)

	breakStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#45B7D1")).
			Bold(true)

	idleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FECA57")).
			Bold(true)

	progressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FECA57")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#95A5A6")).
			Margin(2, 0)
)

type Model struct {
	session *session.Session
	timer   *time.Timer
	running bool
	width   int
	height  int
}

func NewModel(session *session.Session) Model {
	return Model{
		session: session,
		timer:   nil,
		running: false,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "w":
			if !m.running {
				m.session.StartWork()
				m.running = true
				return m, m.startTimer()
			}
			return m, nil
		case "b":
			if !m.running {
				m.session.StartBreak()
				m.running = true
				return m, m.startTimer()
			}
			return m, nil
		case "e":
			m.session.EndSession()
			m.running = false
			if m.timer != nil {
				m.timer.Stop()
			}
			return m, nil
		case "q", "ctrl+c":
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
	title := m.renderTitle()
	status := m.renderStatus()
	progress := m.renderProgress()
	stats := m.renderStats()
	help := m.renderHelp()

	content := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n\n%s",
		title, status, progress, stats, help)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, content)
}

type tickMsg time.Time

func (m Model) startTimer() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) renderTitle() string {
	return titleStyle.Render("POMS ⏳")
}

func (m Model) renderStatus() string {
	var status, statusColor string

	if m.running {
		if m.session.IsWork {
			status = "👷🏾 WORKING..."
			statusColor = workStyle.Render(status)
		} else {
			status = "🍵 BREAK TIME"
			statusColor = breakStyle.Render(status)
		}
	} else {
		status = "💤 IDLE"
		statusColor = idleStyle.Render(status)
	}

	return statusStyle.Render(fmt.Sprintf("Status: %s", statusColor))
}

func (m Model) renderProgress() string {
	if !m.running || m.session.Duration == 0 {
		return ""
	}

	elapsed := time.Since(m.session.StartTime)
	remaining := max(m.session.Duration-elapsed, 0)

	progress := float64(elapsed) / float64(m.session.Duration)
	if progress > 1 {
		progress = 1
	}

	progressBar := m.renderProgressBar(progress)
	timeRemaining := m.renderTimeRemaining(remaining)

	return fmt.Sprintf("%s\n%s", progressBar, timeRemaining)
}

func (m Model) renderProgressBar(progress float64) string {
	barWidth := 30
	filled := int(progress * float64(barWidth))

	return progressStyle.Render(
		fmt.Sprintf("[%s%s] %.0f%%",
			strings.Repeat("█", filled),
			strings.Repeat("░", barWidth-filled),
			progress*100))
}

func (m Model) renderTimeRemaining(remaining time.Duration) string {
	return progressStyle.Render(
		fmt.Sprintf("Time Remaining: %02d:%02d",
			int(remaining.Minutes()),
			int(remaining.Seconds())%60))
}

func (m Model) renderStats() string {
	return fmt.Sprintf(`📊 Stats:
	Sessions Completed: %d
	Current Session: %d`,
		m.session.TotalPomodoros,
		m.session.SessionNumber)
}

func (m Model) renderHelp() string {
	return helpStyle.Render(`🎮 Controls:
	w - Start Work Session
	b - Start Break
	e - End Current Session
	q - Quit`)
}
