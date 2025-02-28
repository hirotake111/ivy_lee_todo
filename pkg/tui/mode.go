package tui

type globalDisplayMode int
type listMode int
type formMode int

const (
	listDisplayMode globalDisplayMode = 10
	formDisplayMode globalDisplayMode = 11

	actionableListMode listMode = 100
	plannedListMode    listMode = 101

	newTaskMode  formMode = 1000
	editTaskMode formMode = 1001
)

type appMode struct {
	globalDisplayMode globalDisplayMode
	currentListMode   listMode
	prevListMode      listMode
	formMode          formMode
}

func newAppMode() appMode {
	return appMode{
		globalDisplayMode: listDisplayMode,
		currentListMode:   actionableListMode,
		prevListMode:      actionableListMode,
		formMode:          newTaskMode,
	}
}

func (a appMode) isListMode() bool {
	return a.globalDisplayMode == listDisplayMode
}

func (a appMode) isNewTaskMode() bool {
	return a.formMode == newTaskMode
}

func (a appMode) isEditTaskMode() bool {
	return a.formMode == editTaskMode
}

func (a appMode) showingActionable() bool {
	return a.currentListMode == actionableListMode
}

func (a *appMode) toggleListMode() *appMode {
	if a.currentListMode == actionableListMode {
		a.currentListMode = plannedListMode
	} else {
		a.currentListMode = actionableListMode
	}
	return a
}

func (a *appMode) toNewTaskMode() *appMode {
	a.prevListMode = a.currentListMode
	a.formMode = newTaskMode
	a.globalDisplayMode = formDisplayMode
	return a
}

func (a *appMode) toEditTaskMode() *appMode {
	a.prevListMode = a.currentListMode
	a.formMode = editTaskMode
	a.globalDisplayMode = formDisplayMode
	return a
}

func (a *appMode) toPrevMode() *appMode {
	a.currentListMode, a.prevListMode = a.prevListMode, a.currentListMode
	a.globalDisplayMode = listDisplayMode
	return a
}
