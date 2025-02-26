package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hirotake111/ivy_lee_todo/pkg/apperrors"
	"github.com/hirotake111/ivy_lee_todo/pkg/domain"
	"github.com/hirotake111/ivy_lee_todo/pkg/service"
)

// General stuff for styling the view
var (
	mainStyle           = lipgloss.NewStyle().MarginLeft(2)
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	docStyle            = lipgloss.NewStyle().Margin(1, 2)

	// Submit Button
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	focusedButton = focusedStyle.Render("[ Submit ]")
)

type model struct {
	ctx             context.Context
	Quitting        bool              // Is the program quitting now?
	TaskList        domain.TaskList   // This includes both planned and actionable tasks. You can get either through receiver methods
	service         *service.Service  // service object
	displayMode     displayMode       // Whether it's showing the list of actionable tasks
	prevDisplayMode displayMode       // Previous display mode
	err             error             // An error message
	cursorMode      cursor.Mode       // Behavior of the cursor (not neccessary?)
	inputIndex      int               // Index of the text input fields
	inputs          []textinput.Model // text input fields
	list            list.Model        // A list of tasks currently been displayed
	showActionable  bool              // True when showing a list of available tasks
	mode            appMode           // app mode
}

func InitializeModel(ctx context.Context, service *service.Service) model {
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
	form.CharLimit = 512
	form.Placeholder = "Description"
	form.PromptStyle = noStyle
	form.TextStyle = noStyle
	forms = append(forms, form)

	list := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	list.Title = "TODO"
	list.SetShowTitle(true)

	return model{
		Quitting:        false,
		ctx:             ctx,
		service:         service,
		displayMode:     actionableListMode,
		prevDisplayMode: actionableListMode,
		inputs:          forms,
		list:            list,
		showActionable:  true,
		mode:            newAppMode(),
	}

}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	return fetchTaskListCmd(m)
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

// selectedTask returns a list of tasks currently been displayed
func (m model) selectedTask() (*domain.Task, error) {
	tasks := m.TaskList.ActionableTasks()
	if !m.mode.showingActionable() {
		tasks = m.TaskList.PlannedTasks()
	}
	if len(tasks) <= m.list.Index() {
		return nil, apperrors.OutOfIndex{}
	}
	return tasks[m.list.Index()], nil
}

// Update implements tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.err = nil // Refresh error message
	if m.mode.isListMode() {
		return updateWithActionableListMode(m, msg)
	} else {
		return updateWithTaskEditMode(m, msg)
	}
}

func updateWithActionableListMode(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	// Triggered by a key stroke
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		// Toggle actionable list mode to planned list mode
		case "t":
			m.mode.toggleListMode()
			if m.mode.showingActionable() {
				m.list.Title = "TODO"
			} else {
				m.list.Title = "Planned"
			}
			return m, refleshDisplayItemCmd(m)

		// Move selected task into planned one
		case "m":
			task, err := m.selectedTask()
			if err != nil {
				break
			}
			if m.mode.showingActionable() {
				return m, updateTaskCmd(m, task.ToPlanned())
			}
			return m, updateTaskCmd(m, task.ToActionable())

		// Complete selected task
		case " ":
			selected, err := m.selectedTask()
			if err != nil {
				break
			}
			// Only performs completion when it's displaying actionable tasks
			if m.mode.showingActionable() {
				return m, completeTaskCmd(m, selected)
			}

		// Delete a selected task
		case "d":
			selected, err := m.selectedTask()
			if err != nil {
				break
			}
			return m, deleteTaskCmd(m, selected)

		// New task mode
		case "a":
			// Update the view only when user is not filtering
			if !(m.list.FilterState() != list.Unfiltered) {
				m.mode.toNewTaskMode()
			}
		}

	// Triggered by a new task list
	case listMsg:
		// Update task list
		m.TaskList = msg.tasks
		//Reflesh items currently been displayed
		return m, refleshDisplayItemCmd(m)

	// Reflesh display items
	case newDisplayItemsMsg:
		m.list.ResetSelected()
		return m, m.list.SetItems(NewDisplayItems(msg.tasks))

	// Triggered by error message
	case errMsg:
		m.err = msg
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func updateWithTaskEditMode(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	// Triggered by a key stroke
	case tea.KeyMsg:
		switch msg.Type {
		// Back to the previous mode
		case tea.KeyCtrlC, tea.KeyEscape:
			// m.displayMode = m.prevDisplayMode
			m.mode.toPrevMode()
			m.ClearInputs()
			return m, nil

		// Go to the next input
		case tea.KeyTab, tea.KeyEnter:
			// Perform submission only when the button is focused & end user hits the enter key.
			if msg.Type == tea.KeyEnter && m.canSubmit() {
				m.mode.toPrevMode()
				req := &domain.NewTaskRequest{
					Title:       m.inputs[0].Value(),
					Description: m.inputs[1].Value(),
				}
				m.ClearInputs()
				return m, newTaskCmd(m, req)
			}
			var cmd tea.Cmd
			m.inputIndex = (m.inputIndex + 1) % (len(m.inputs) + 1)
			m.updateInputStyle()
			return m, cmd
		}
	}

	// Handle character input and blinking
	return m, m.updateInput(msg)
}

func (m *model) updateInput(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

// View implements tea.Model.
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}

	if m.mode.isListMode() {
		return docStyle.Render(m.list.View())
	}

	// New task form
	s = newTaskFormView(m)

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
