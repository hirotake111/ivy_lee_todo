package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/hirotake111/ivy_lee_todo/pkg/domain"
)

// displayItem is a container for bubbles/list item
type displayItem struct {
	title, desc string
}

func NewDisplayItems(tasks []*domain.Task) (is []list.Item) {
	for _, t := range tasks {
		is = append(is, displayItem{title: t.Title(), desc: t.Description()})
	}
	return
}
func (i displayItem) Title() string       { return i.title }
func (i displayItem) Description() string { return i.desc }
func (i displayItem) FilterValue() string { return i.title }
