package tui

import (
	"sort"

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
		sort.Slice(l, func(i, j int) bool {
			return l[i].CreatedAt().After(l[j].CreatedAt())
		})
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
			return errMsg{err: err}
		}
		return tasksUpdatedMsg{}
	}
}

// updateTaskCmd updates a task and fetches the latest task list
func updateTaskCmd(m model, task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.Update(m.ctx, task); err != nil {
			return errMsg{err: err}
		}
		return tasksUpdatedMsg{}
	}
}

// makeActionableTaskCmd makes a task actionable
func makeActionableTaskCmd(m model, task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.MakeActionable(m.ctx, task.Id()); err != nil {
			return errMsg{err: err}
		}
		return tasksUpdatedMsg{}
	}
}

// completeTaskCmd updates a task and fetches the latest task list
func completeTaskCmd(m model, task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.DeleteTask(m.ctx, task.Id()); err != nil {
			return errMsg{err: err}
		}
		return tasksUpdatedMsg{}
	}
}

// deleteTaskCmd deletes a task and fetches the latest task list
func deleteTaskCmd(m model, task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.DeleteTask(m.ctx, task.Id()); err != nil {
			return errMsg{err: err}
		}
		return tasksUpdatedMsg{}
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

// successCmd generates successMsg using passed text value
func successCmd(text string) tea.Cmd {
	return func() tea.Msg {
		return successMsg{text: text}
	}
}

// goToEditModeCmd generates toEditModeMsg
func goToEditModeCmd(m model) tea.Cmd {
	return func() tea.Msg {
		t, err := m.selectedTask()
		if err != nil {
			return errMsg{err: err}
		}
		return toEditModeMsg{task: t}
	}
}
