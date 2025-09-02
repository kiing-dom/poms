package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kiing-dom/poms/internal/audio"
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

	state       string
	configIndex int
	textInputs  []textinput.Model
	presets     []TimePreset
}

var defaultPresets = []TimePreset{
	{"Classic", 25, 5},
	{"Short Focus", 15, 3},
	{"Deep Work", 45, 15},
	{"Flow State", 90, 20},
	{"Custom", 0, 0}, // allow manual input
}

type TimePreset struct {
	Name         string
	WorkMinutes  int
	BreakMinutes int
}

func NewModel(session *session.Session) Model {
	workInput := textinput.New()
	workInput.Placeholder = "25"
	workInput.CharLimit = 3
	workInput.Focus()

	breakInput := textinput.New()
	breakInput.Placeholder = "5"
	breakInput.CharLimit = 3

	return Model{
		session:     session,
		timer:       nil,
		running:     false,
		width:       0, // Start with 0 to force detection
		height:      0, // Start with 0 to force detection
		state:       "config",
		configIndex: 0,
		textInputs:  []textinput.Model{workInput, breakInput},
		presets:     defaultPresets,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		func() tea.Msg {
			return tea.WindowSizeMsg{Width: 80, Height: 24} // Fallback size
		},
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Ensure we have minimum viable dimensions
		if m.width < 10 {
			m.width = 10
		}
		if m.height < 5 {
			m.height = 5
		}
		return m, nil	
	case tea.KeyMsg:
		switch m.state {
		case "config":
			return m.updateConfig(msg)
		case "timer":
			return m.updateTimer(msg)
		}
	case tickMsg:
		if m.running && m.session.IsSessionActive() {
			return m, m.startTimer()
		}

		if m.running && m.session.IsWork {
			audio.PlayNotification("assets/sounds/timer-beep.mp3")
			m.session.StartBreak()
			m.running = true
			return m, m.startTimer()
		}
		m.running = false
		return m, nil
	}

	return m, nil
}

func (m *Model) updateConfig(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.configIndex > 0 {
			m.configIndex--
		}
	case "down", "j":
		if m.configIndex < len(m.presets)-1 {
			m.configIndex++
		}
	case "enter":
		return m.applyConfiguration()
	case "tab":
		if m.presets[m.configIndex].Name == "Custom" {
			if m.textInputs[0].Focused() {
				m.textInputs[0].Blur()
				m.textInputs[1].Focus()
			} else {
				m.textInputs[1].Blur()
				m.textInputs[0].Focus()
			}
		}
	case "q", "ctrl + c":
		return m, tea.Quit
	default:
		if m.presets[m.configIndex].Name == "Custom" {
			var cmd tea.Cmd
			if m.textInputs[0].Focused() {
				m.textInputs[0], cmd = m.textInputs[0].Update(msg)
			} else if m.textInputs[1].Focused() {
				m.textInputs[1], cmd = m.textInputs[1].Update(msg)
			}

			return m, cmd
		}
	}
	return m, nil
}

func (m *Model) updateTimer(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "c":
		m.state = "config"
		return m, nil
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
	return m, nil
}

func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		m.width = 80
		m.height = 24
	}

	switch m.state {
	case "config":
		return m.renderConfig()
	case "timer":
		return m.renderTimer()
	default:
		return m.renderConfig()
	}
}

func (m Model) getLayoutType() string {
	switch {
	case m.width < 30 || m.height < 8:
		return "minimal"
	case m.width < 50 || m.height < 12:
		return "compact"
	case m.width < 80 || m.height < 18:
		return "medium"
	default:
		return "full"
	}
}

type tickMsg time.Time

func (m Model) startTimer() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) renderTitle() string {
	layout := m.getLayoutType()

	switch layout {
	case "minimal":
		return "" // No title in minimal view
	case "compact":
		return titleStyle.Render("POMS")
	default:
		return titleStyle.Render("POMS ⏳")
	}
}

func (m Model) renderStatus() string {
	var status, statusColor string
	layout := m.getLayoutType()

	if m.running {
		if m.session.IsWork {
			if layout == "minimal" || layout == "compact" {
				status = "W"
			} else {
				status = "WORKING..."
			}
			statusColor = workStyle.Render(status)
		} else {
			if layout == "minimal" || layout == "compact" {
				status = "B"
			} else {
				status = "BREAK"
			}
			statusColor = breakStyle.Render(status)
		}
	} else {
		if layout == "minimal" || layout == "compact" {
			status = "I"
		} else {
			status = "IDLE"
		}
		statusColor = idleStyle.Render(status)
	}

	if layout == "minimal" {
		return statusColor
	}

	return statusStyle.Render(fmt.Sprintf("Status: %s", statusColor))
}

func (m Model) renderProgress() string {
	if !m.running || (m.session.IsWork && m.session.WorkDuration == 0) || (!m.session.IsWork && m.session.BreakDuration == 0) {
		return ""
	}

	elapsed := time.Since(m.session.StartTime)
	currentDuration := m.session.GetCurrentDuration()
	remaining := max(currentDuration-elapsed, 0)

	progress := float64(elapsed) / float64(currentDuration)
	if progress > 1 {
		progress = 1
	}

	progressBar := m.renderProgressBar(progress)
	timeRemaining := m.renderTimeRemaining(remaining)

	return fmt.Sprintf("%s\n%s", progressBar, timeRemaining)
}

func (m Model) renderProgressBar(progress float64) string {
	// Make progress bar width responsive with safe minimums
	layout := m.getLayoutType()
	var barWidth int

	switch layout {
	case "minimal":
		barWidth = max(5, min(10, m.width-5))
	case "compact":
		barWidth = max(8, min(20, m.width-8))
	case "medium":
		barWidth = max(15, min(35, m.width-10))
	default:
		barWidth = max(20, min(50, m.width-15))
	}

	// Ensure we don't have negative or zero width
	if barWidth <= 0 {
		barWidth = 10
	}

	filled := int(progress * float64(barWidth))
	if filled < 0 {
		filled = 0
	}
	if filled > barWidth {
		filled = barWidth
	}

	if layout == "minimal" {
		// Ultra compact version
		return progressStyle.Render(
			fmt.Sprintf("[%s%s]",
				strings.Repeat("█", filled),
				strings.Repeat("░", barWidth-filled)))
	}

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
	layout := m.getLayoutType()

	switch layout {
	case "minimal":
		return "" // No stats in minimal view
	case "compact":
		if m.width >= 20 {
			return fmt.Sprintf("📊 %d/%d", m.session.TotalPomodoros, m.session.SessionNumber)
		}
		return ""
	case "medium":
		if m.width >= 35 {
			return fmt.Sprintf("📊 Current: %d | Sessions: %d",
				m.session.SessionNumber, m.session.TotalPomodoros)
		}
		return ""
	default:
		return fmt.Sprintf(`📊 Stats:
		Sessions Completed: %d
		Current Session: %d`,
			m.session.TotalPomodoros,
			m.session.SessionNumber)
	}
}

func (m Model) renderHelp() string {
	layout := m.getLayoutType()

	switch layout {
	case "minimal":
		return ""
	case "compact":
		if m.width >= 25 {
			return helpStyle.Render("w:Work b:Break c:Config e:End q:Quit")
		}
		return helpStyle.Render("w/b/c/e/q")
	case "medium":
		if m.width >= 45 {
			return helpStyle.Render("Controls: w-Work b-Break c-Config e-End q-Quit")
		}
		return helpStyle.Render("w:Work b:Break c:Config e:End q:Quit")
	default:
		return helpStyle.Render(`🎮 Controls:
		w - Start Work Session
		b - Start Break
		c - Configure Times
		e - End Current Session
		q - Quit`)
	}
}

func (m Model) renderConfig() string {
	layout := m.getLayoutType()

	var sections []string

	title := titleStyle.Render("POMS CONFIG")
	sections = append(sections, title)

	presetSection := m.renderPresetOptions()
	sections = append(sections, presetSection)

	if m.presets[m.configIndex].Name == "Custom" {
		customSection := m.renderCustomInputs()
		sections = append(sections, customSection)
	}

	helpText := m.renderConfigHelp()
	sections = append(sections, helpText)

	content := strings.Join(sections, "\n\n")

	if layout != "minimal" {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

func (m Model) renderPresetOptions() string {
	var options []string

	for i, preset := range m.presets {
		var option string
		if preset.Name == "Custom" {
			option = preset.Name
		} else {
			option = fmt.Sprintf("%s (%dm work, %dm break)", preset.Name, preset.WorkMinutes, preset.BreakMinutes)
		}

		if i == m.configIndex {
			option = workStyle.Render("> " + option)
		} else {
			option = "  " + option
		}
		options = append(options, option)
	}

	return strings.Join(options, "\n")
}

func (m Model) renderCustomInputs() string {
	workLabel := "Work Duration (minutes):"
	breakLabel := "Break Duration (minutes):"

	workLine := workLabel + m.textInputs[0].View()
	breakLine := breakLabel + m.textInputs[1].View()

	return fmt.Sprintf("%s\n%s", workLine, breakLine)
}

func (m Model) renderConfigHelp() string {
	layout := m.getLayoutType()

	switch layout {
	case "minimal":
		return helpStyle.Render("↑↓:Select Enter:Confirm Tab:Switch q:Quit")
	case "compact":
		return helpStyle.Render("↑↓: Select | Enter: Confirm | Tab: Switch | q: Quit")
	default:
		return helpStyle.Render(`Configuration Controls:
		↑/↓ or j/k - Navigate presets
		Enter - Apply configuration and start timer
		Tab - Switch between custom input fields
		q - Quit`)
	}
}

func (m Model) renderTimer() string {
	// Determine layout based on actual terminal size
	layout := m.getLayoutType()

	var sections []string

	// Build sections based on layout
	title := m.renderTitle()
	if title != "" {
		sections = append(sections, title)
	}

	status := m.renderStatus()
	if status != "" {
		sections = append(sections, status)
	}

	progress := m.renderProgress()
	if progress != "" {
		sections = append(sections, progress)
	}

	stats := m.renderStats()
	if stats != "" {
		sections = append(sections, stats)
	}

	help := m.renderHelp()
	if help != "" {
		sections = append(sections, help)
	}

	// Join sections with appropriate spacing based on layout
	var content string
	switch layout {
	case "minimal":
		content = lipgloss.JoinVertical(lipgloss.Left, sections...)
	case "compact":
		content = lipgloss.JoinVertical(lipgloss.Left, sections...)
	default:
		// Add more spacing for larger screens
		content = strings.Join(sections, "\n\n")
	}

	// Only use Place if we have reliable dimensions and enough space
	contentWidth := lipgloss.Width(content)
	contentHeight := lipgloss.Height(content)

	if m.width > 0 && m.height > 0 &&
		contentWidth <= m.width && contentHeight <= m.height &&
		layout != "minimal" {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}

func (m *Model) applyConfiguration() (tea.Model, tea.Cmd) {
	preset := m.presets[m.configIndex]

	var workMinutes, breakMinutes int
	var err error

	if preset.Name == "Custom" {
		workStr := m.textInputs[0].Value()
		breakStr := m.textInputs[1].Value()

		if workStr == "" {
			workStr = m.textInputs[0].Placeholder
		}
		if breakStr == "" {
			breakStr = m.textInputs[1].Placeholder
		}

		workMinutes, err = strconv.Atoi(workStr)
		if err != nil || workMinutes <= 0 {
			return m, nil
		}

		breakMinutes, err = strconv.Atoi(breakStr)
		if err != nil || breakMinutes <= 0 {
			return m, nil
		}
	} else {
		workMinutes = preset.WorkMinutes
		breakMinutes = preset.BreakMinutes
	}

	m.session.WorkDuration = time.Duration(workMinutes) * time.Minute
	m.session.BreakDuration = time.Duration(breakMinutes) * time.Minute

	m.state = "timer"
	return m, nil
}
