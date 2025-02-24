package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hirotake111/ivy_lee_todo/pkg/domain"
)

// fetchTaskListCmd generates a command that sends ListMsg
func fetchTaskListCmd(m model) tea.Cmd {
	return func() tea.Msg {
		l, err := m.service.List(m.ctx)
		if err != nil {
			return errMsg{err: err}
		}
		return listMsg{tasks: l}
	}
}

// errCmd generates a command that sends errMsg
func errCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errMsg{err: err}
	}
}

// newTaskCmd generates a command that creates a new planned task
func newTaskCmd(m model, task *domain.NewTaskRequest) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.AddTask(m.ctx, task.Title, task.Description); err != nil {
			return errCmd(err)
		}
		return fetchTaskListCmd(m)()
	}
}

// updateTaskCmd updates a task and fetches the latest task list
func updateTaskCmd(m model, task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.Update(m.ctx, task); err != nil {
			return errMsg{err: err}
		}
		return fetchTaskListCmd(m)()
	}
}

// completeTaskCmd updates a task and fetches the latest task list
func completeTaskCmd(m model, task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.DeleteTask(m.ctx, task.Id()); err != nil {
			return errMsg{err: err}
		}
		return fetchTaskListCmd(m)()
	}
}

// deleteTaskCmd deletes a task and fetches the latest task list
func deleteTaskCmd(m model, task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.DeleteTask(m.ctx, task.Id()); err != nil {
			return errMsg{err: err}
		}
		return fetchTaskListCmd(m)()
	}
}

// refleshDisplayItemCmd generates a command that sends a message with new display items
func refleshDisplayItemCmd(m model) tea.Cmd {
	return func() tea.Msg {
		if m.mode.showingActionable() {
			return newDisplayItemsMsg{tasks: m.TaskList.ActionableTasks()}
		}
		return newDisplayItemsMsg{tasks: m.TaskList.PlannedTasks()}
	}
}
