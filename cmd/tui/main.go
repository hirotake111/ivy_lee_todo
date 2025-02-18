package main

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fogleman/ease"
	"github.com/hirotake111/ivy_lee_todo/pkg/db"
	"github.com/hirotake111/ivy_lee_todo/pkg/domain"
	"github.com/hirotake111/ivy_lee_todo/pkg/repository"
	"github.com/hirotake111/ivy_lee_todo/pkg/service"
	"github.com/lucasb-eyer/go-colorful"
)

type displayMode int

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "

	actionableListMode displayMode = iota
	plannedListMode
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
	ctx := context.Background()
	db := db.NewSqlite3Db(false)
	r := repository.NewSQLiteRepository()
	s := service.NewService(db, r)
	m := initializeModel(ctx, s)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Couldn't start program: %s\n", err)
	}
}

type model struct {
	ctx         context.Context
	Choice      int              // Index of the item cursor is currently pointing to
	Chosen      bool             // Whether the item is chosen or not
	Ticks       int              // Tick!
	Frames      int              // Frames?
	Progress    float64          // not neccessary
	Loaded      bool             // Whether it's loaded
	Quitting    bool             // Is the program quitting now?
	TaskList    domain.TaskList  // This includes both planned and actionable tasks. You can get either through receiver methods
	service     *service.Service // service object
	displayMode displayMode      // Whether it's showing the list of actionable tasks
}

func initializeModel(ctx context.Context, service *service.Service) model {
	return model{
		Choice:      0,
		Chosen:      false,
		Ticks:       10,
		Frames:      0,
		Progress:    0,
		Loaded:      false,
		Quitting:    false,
		ctx:         ctx,
		service:     service,
		displayMode: actionableListMode,
	}

}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	return m.fetchTaskList
}

// fetchTaskList retrieves a list of tasks and returns ListMsg
func (m model) fetchTaskList() tea.Msg {
	l, err := m.service.List(m.ctx)
	if err != nil {
		panic(err)
	}
	return ListMsg{tasks: l}
}

// updateTask updates a task and fetches the latest task list
func (m model) updateTask(task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.Update(m.ctx, task); err != nil {
			return errMsg{err: err}
		}
		return m.fetchTaskList()
	}
}

// Update implements tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.displayMode {

	// actionable list displayMode
	case actionableListMode:
		switch msg := msg.(type) {
		// Triggered by a key stroke
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				m.Quitting = true
				return m, tea.Quit

			// Switch actionable list mode to planned list mode
			case "r":
				m.Choice = 0 // Reset choice index
				m.displayMode = plannedListMode
				return m, nil

			case "m":
				m.Choice = 0 // Reset choice index
				t := m.TaskList.ActionableTasks()[m.Choice]
				return m, m.updateTask(t.ToPlanned())
			}

		// Triggered by a new task list
		case ListMsg:
			m.TaskList = msg.tasks
			return m, nil
		}

		// Hand off the message and model to the appropriate update function for the
		// appropriate view based on the current state.
		if !m.Chosen {
			return updateChoices(msg, m)
		}

		// planned list displayMode
	case plannedListMode:
		switch msg := msg.(type) {

		// Triggered by a key stroke
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				m.Quitting = true
				return m, tea.Quit

			// Switch planed list mode to actionable list mode
			case "r":
				m.Choice = 0 // Reset choice index
				m.displayMode = actionableListMode
				return m, nil

			// Make planned task into actionable one
			case " ":
				t := m.TaskList.PlannedTasks()[m.Choice]
				m.Choice = 0 // Resent choice index
				m.displayMode = actionableListMode
				return m, m.updateTask(t.ToActionable())
			}

		// Triggered by a new task list
		case ListMsg:
			m.TaskList = msg.tasks
			return m, nil
		}

		// Hand off the message and model to the appropriate update function for the
		// appropriate view based on the current state.
		if !m.Chosen {
			return updateChoices(msg, m)
		}
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
	var tasks []*domain.Task
	if m.displayMode == actionableListMode {
		tasks = m.TaskList.ActionableTasks()
	} else if m.displayMode == plannedListMode {
		tasks = m.TaskList.PlannedTasks()
	}
	var tpl string
	if m.displayMode == actionableListMode {
		tpl = fmt.Sprintf("TODOs (%d/%d)\n\n", len(tasks), len(m.TaskList))
	} else if m.displayMode == plannedListMode {
		tpl = fmt.Sprintf("Planned Tasks (%d/%d)\n\n", len(tasks), len(m.TaskList))
	}
	tpl += "%s\n\n"
	tpl += helpMessage(&m)
	var choices strings.Builder
	for i, t := range tasks {
		cb := checkbox(t.Title(), c == i)
		choices.WriteString(fmt.Sprintf("%s\n", cb))
	}
	return fmt.Sprintf(tpl, choices.String())
}

func helpMessage(m *model) string {
	vertical := subtleStyle.Render("j/k, up/down: select") + dotStyle
	edit := subtleStyle.Render("e: edit") + dotStyle
	actionable := subtleStyle.Render("m: make it actionable") + dotStyle
	done := subtleStyle.Render("d: done") + dotStyle
	del := subtleStyle.Render("d: delete") + dotStyle
	quit := subtleStyle.Render("q, esc: quit")
	if m.displayMode == actionableListMode {
		return vertical + edit + done + quit
	}
	return vertical + actionable + edit + del + quit
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] -" + label)
	}
	return fmt.Sprintf("[ ] -%s", label)
}

type (
	tickMsg  struct{}
	frameMsg struct{}
)

type ListMsg struct {
	tasks domain.TaskList
}

// errMsg is a message that contains an error in it.
type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

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
	l := len(m.TaskList.ActionableTasks())
	if m.displayMode == plannedListMode {
		l = len(m.TaskList.PlannedTasks())
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice = min(m.Choice+1, l-1)
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

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
