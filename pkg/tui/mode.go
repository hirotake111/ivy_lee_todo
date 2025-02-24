package tui

type appMode struct {
	mode displayMode
	prev displayMode
}

func newAppMode() appMode {
	return appMode{
		mode: actionableListMode,
		prev: actionableListMode,
	}
}

func (a appMode) isListMode() bool {
	return a.mode == actionableListMode || a.mode == plannedListMode
}

func (a appMode) showingActionable() bool {
	return a.mode == actionableListMode
}

func (a *appMode) toggleListMode() *appMode {
	if a.mode == actionableListMode {
		a.mode = plannedListMode
	} else if a.mode == plannedListMode {
		a.mode = actionableListMode
	}
	return a
}

func (a *appMode) toNewTaskMode() *appMode {
	a.prev = a.mode
	a.mode = newTaskMode
	return a
}

func (a *appMode) toPrevMode() *appMode {
	a.mode, a.prev = a.prev, a.mode
	return a
}
