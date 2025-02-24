package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hirotake111/ivy_lee_todo/pkg/apperrors"
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
	newTaskMode
	editTaskMode
)

// General stuff for styling the view
var (
	keywordStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty       = subtleStyle.Render(progressEmptyChar)
	dotStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle           = lipgloss.NewStyle().MarginLeft(2)
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()

	// Gradient colors we'll use for the progress bar
	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
	// Submit Button
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	focusedButton = focusedStyle.Render("[ Submit ]")
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
	ctx             context.Context
	Choice          int               // Index of the item cursor is currently pointing to
	Ticks           int               // Tick!
	Frames          int               // Frames?
	Progress        float64           // not neccessary
	Loaded          bool              // Whether it's loaded
	Quitting        bool              // Is the program quitting now?
	TaskList        domain.TaskList   // This includes both planned and actionable tasks. You can get either through receiver methods
	service         *service.Service  // service object
	displayMode     displayMode       // Whether it's showing the list of actionable tasks
	prevDisplayMode displayMode       // Previous display mode
	err             error             // An error message
	cursorMode      cursor.Mode       // Behavior of the cursor (not neccessary?)
	inputIndex      int               // Index of the text input fields
	inputs          []textinput.Model // text input fields
}

func initializeModel(ctx context.Context, service *service.Service) model {
	forms := make([]textinput.Model, 0)
	var form textinput.Model
	// Title
	form = textinput.New()
	form.Cursor.Style = cursorStyle
	form.CharLimit = 32
	form.Placeholder = "Title"
	form.Focus()
	form.PromptStyle = focusedStyle
	form.TextStyle = focusedStyle
	forms = append(forms, form)
	// Description
	form = textinput.New()
	form.Cursor.Style = cursorStyle
	form.CharLimit = 32
	form.Placeholder = "Description"
	form.PromptStyle = noStyle
	form.TextStyle = noStyle
	forms = append(forms, form)

	return model{
		Choice:          0,
		Ticks:           10,
		Frames:          0,
		Progress:        0,
		Loaded:          false,
		Quitting:        false,
		ctx:             ctx,
		service:         service,
		displayMode:     actionableListMode,
		prevDisplayMode: actionableListMode,
		inputs:          forms,
	}

}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	return fetchTaskListCmd(&m)
}

// updateInputStyle recalculates each input style based on focused index
func (m *model) updateInputStyle() {
	for i := 0; i < len(m.inputs); i++ {
		if i == m.inputIndex {
			// Set focused state
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else {
			// Remove focused state
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
	}
}

func (m model) canSubmit() bool {
	return m.inputIndex == len(m.inputs) && len(m.inputs[0].Value()) > 0
}

// ClearInputs clears all inputs
func (m *model) ClearInputs() {
	for i := range m.inputs {
		m.inputs[i].Reset()
	}
	// Update input style based on focused index
	m.inputIndex = 0
	m.updateInputStyle()
}

// fetchTaskListCmd retrieves a list of tasks and returns ListMsg
func fetchTaskListCmd(m *model) tea.Cmd {
	return func() tea.Msg {
		l, err := m.service.List(m.ctx)
		if err != nil {
			panic(err)
		}
		return ListMsg{tasks: l}
	}
}

// errCmd generates a command that sends errMsg
func errCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errMsg{err: err}
	}
}

// newTaskCmd generates a command that creates a new planned task
func newTaskCmd(m *model, task *domain.NewTaskRequest) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.AddTask(m.ctx, task.Title, task.Description); err != nil {
			return errCmd(err)
		}
		return fetchTaskListCmd(m)()
	}
}

// availableTaskCmd generates a command that turns a task into actionable.
// If the num of tasks exceeds the limit, it sends errMsg.
func availableTaskCmd(m *model, task *domain.Task) tea.Cmd {
	if !m.TaskList.CanAddAnother() {
		return errCmd(apperrors.NewTaskExceededError(m.TaskList.MaxTskNum()))
	}
	return updateTaskCmd(m, task.ToActionable())
}

// updateTaskCmd updates a task and fetches the latest task list
func updateTaskCmd(m *model, task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.Update(m.ctx, task); err != nil {
			return errMsg{err: err}
		}
		return fetchTaskListCmd(m)()
	}
}

// completeTaskCmd updates a task and fetches the latest task list
func completeTaskCmd(m *model, task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.DeleteTask(m.ctx, task.Id()); err != nil {
			return errMsg{err: err}
		}
		return fetchTaskListCmd(m)()
	}
}

// deleteTaskCmd deletes a task and fetches the latest task list
func deleteTaskCmd(m *model, task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.DeleteTask(m.ctx, task.Id()); err != nil {
			return errMsg{err: err}
		}
		return fetchTaskListCmd(m)()
	}
}

// Update implements tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.err = nil // Refresh error message
	switch m.displayMode {

	// actionable list displayMode
	case actionableListMode:
		return updateWithActionableListMode(m, msg)

	// planned list displayMode
	case plannedListMode:
		return updateWithPlannedlistMode(m, msg)

	// edit task mode
	case newTaskMode, editTaskMode:
		switch msg := msg.(type) {
		// Triggered by a key stroke
		case tea.KeyMsg:
			switch msg.Type {
			// Back to the previous mode
			case tea.KeyCtrlC, tea.KeyEscape:
				m.displayMode = m.prevDisplayMode
				m.ClearInputs()
				return m, nil

			// Go to the next input
			case tea.KeyTab, tea.KeyEnter:
				// If submit button is focused end user hit enter, then submit the new item
				if msg.Type == tea.KeyEnter && m.canSubmit() {
					m.displayMode = m.prevDisplayMode
					req := &domain.NewTaskRequest{
						Title:       m.inputs[0].Value(),
						Description: m.inputs[1].Value(),
					}
					m.ClearInputs()
					return m, newTaskCmd(&m, req)
				}
				var cmd tea.Cmd
				m.inputIndex = (m.inputIndex + 1) % (len(m.inputs) + 1)
				// Update input style based on focused index
				m.updateInputStyle()
				return m, cmd
			}
		}

		// Handle character input and blinking
		return m.updateInput(msg)

	default:
		panic(fmt.Sprintf("unknown display mode: %d", m.displayMode))
	}
}

func (m *model) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

// View implements tea.Model.
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}

	if m.displayMode == newTaskMode {
		s = newTaskFormView(m)
	} else {
		s = choicesView(m)
	}
	var errMessage string
	if m.err != nil {
		errMessage = errorStyle.Render(fmt.Sprintf("\nERROR: %s\n", m.err.Error()))
	}
	return mainStyle.Render("\n" + s + errMessage + "\n\n")
}

// newTaskFormView returns form for a new task
func newTaskFormView(m model) string {
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.inputIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
	return b.String()
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
		tpl = fmt.Sprintf("TODOs (%d/%d)\n\n", len(tasks), m.TaskList.MaxTskNum())
	} else if m.displayMode == plannedListMode {
		tpl = fmt.Sprintf("Planned Tasks (%d)\n\n", len(tasks))
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
	switchMode := subtleStyle.Render("s: switch mode") + dotStyle
	edit := subtleStyle.Render("e: edit") + dotStyle
	actionable := subtleStyle.Render("space: make it actionable") + dotStyle
	done := subtleStyle.Render("space: done") + dotStyle
	del := subtleStyle.Render("d: delete") + dotStyle
	quit := subtleStyle.Render("q, esc: quit")
	if m.displayMode == actionableListMode {
		return vertical + done + edit + switchMode + quit
	}
	return vertical + actionable + edit + del + switchMode + quit
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
		}
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

func updateWithActionableListMode(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Triggered by a key stroke
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.Quitting = true
			return m, tea.Quit

		// Toggle actionable list mode to planned list mode
		case "t":
			m.Choice = 0 // Reset choice index
			m.displayMode = plannedListMode
			return m, nil

		// Move selected task into planned one
		case "m":
			m.Choice = 0 // Reset choice index
			t := m.TaskList.ActionableTasks()[m.Choice]
			return m, updateTaskCmd(&m, t.ToPlanned())

		// Complete selected task
		case " ":
			t := m.TaskList.ActionableTasks()[m.Choice]
			return m, completeTaskCmd(&m, t)

		// New task mode
		case "a":
			m.prevDisplayMode = m.displayMode
			m.displayMode = newTaskMode
			return m, nil
		}

	// Triggered by a new task list
	case ListMsg:
		m.TaskList = msg.tasks
		return m, nil

	// Triggered by error message
	case errMsg:
		m.err = msg
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	return updateChoices(msg, m)
}

func updateWithPlannedlistMode(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Triggered by a key stroke
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.Quitting = true
			return m, tea.Quit

		// Toggle planed list mode to actionable list mode
		case "t":
			m.Choice = 0 // Reset choice index
			m.displayMode = actionableListMode
			return m, nil

		// Move selected planned task into actionable
		case "m":
			t := m.TaskList.PlannedTasks()[m.Choice]
			return m, availableTaskCmd(&m, t)

		// New task mode
		case "a":
			m.prevDisplayMode = m.displayMode
			m.displayMode = newTaskMode
			return m, nil

		// Delete a selected task
		case "d":
			t := m.TaskList.PlannedTasks()[m.Choice]
			return m, deleteTaskCmd(&m, t)
		}

	// Triggered by a new task list
	case ListMsg:
		m.TaskList = msg.tasks
		return m, nil

	// Triggered by error message
	case errMsg:
		m.err = msg
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	return updateChoices(msg, m)
}
