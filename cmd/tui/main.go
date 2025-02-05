package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fogleman/ease"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

// General stuff for styling the view
var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)

	// Gradient colors we'll use for the progress bar
	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)

func main() {
	initialModel := model{
		Choice:   0,
		Chosen:   false,
		Ticks:    10,
		Frames:   0,
		Progress: 0,
		Loaded:   false,
		Quitting: false,
	}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Couldn't start program: %s\n", err)
	}
}

type model struct {
	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	return tick()
}

// Update implements tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
}

// View implements tea.Model.
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return mainStyle.Render("\n" + s + "\n\n")
}

// The second view, after a task has beeen chosen
func chosenView(m model) string {
	var msg string
	switch m.Choice {
	case 0:
		msg = fmt.Sprintf("Carrot planting?\n\nCool, we'll need %s and %s...", keywordStyle.Render("libgarden"), keywordStyle.Render("vegeutils"))
	case 1:
		msg = fmt.Sprintf("A trip to the market?\n\nOkay, then we should install %s and %s...", keywordStyle.Render("marketkit"), keywordStyle.Render("libshopping"))
	case 2:
		msg = fmt.Sprintf("Reading time?\n\nOkay, cool, then we’ll need a library. Yes, an %s.", keywordStyle.Render("actual library"))
	default:
		msg = fmt.Sprintf("It’s always good to see friends.\n\nFetching %s and %s...", keywordStyle.Render("social-skills"), keywordStyle.Render("conversationutils"))
	}
	label := "Downloading..."
	if m.Loaded {
		label = fmt.Sprintf("Downloaded. Exisitng in %s seconds.,..", ticksStyle.Render(strconv.Itoa(m.Ticks)))
	}
	return msg + "\n\n" + label + "\n" + progressbar(m.Progress) + "%"
}

func progressbar(percent float64) string {
	w := float64(progressBarWidth)
	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += ramp[i].Render(progressFullChar)
	}
	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)
	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// The first view, where you're choosing a task
func choicesView(m model) string {
	c := m.Choice
	tpl := "What to do today?\n\n"
	tpl += "%s\n\n"
	tpl += "Program quits in %s seconds\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")
	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		checkbox("Plant carrots", c == 0),
		checkbox("Go to the market", c == 1),
		checkbox("Read something", c == 2),
		checkbox("See friends", c == 3),
	)
	return fmt.Sprintf(tpl, choices, ticksStyle.Render(strconv.Itoa(m.Ticks)))
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

type (
	tickMsg  struct{}
	frameMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return frameMsg{}
	})
}

func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	ca, err := colorful.Hex(colorA)
	if err != nil {
		panic(err)
	}
	cb, err := colorful.Hex(colorB)
	if err != nil {
		panic(err)
	}
	for i := 0.0; i < steps; i++ {
		c := ca.BlendLuv(cb, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return s
}

func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

// Update loop for the second view after a choice has beeen made
func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case frameMsg:
		if !m.Loaded {
			m.Frames++
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
			if m.Progress >= 1 {
				m.Progress = 1
				m.Loaded = true
				m.Ticks = 3
				return m, tick()
			}
			return m, frame()
		}
	case tickMsg:
		if m.Loaded {
			if m.Ticks == 0 {
				m.Quitting = true
				return m, tea.Quit
			}
			m.Ticks--
			return m, tick()
		}
	}
	return m, nil
}

// Update loop for the first view where you're choosing a task
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice = min(m.Choice+1, 3)
		case "k", "up":
			m.Choice = max(m.Choice-1, 0)
		case "enter":
			m.Chosen = true
			return m, frame()
		}
	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks--
		return m, tick()
	}
	return m, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
