package taskView

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	"github.com/thoas/go-funk"
)

type TaskViewModel struct {
	// HideSelfWhenDone              bool
	DisplayOnlyChildrenInProgress bool
}

func ViewTask(m task.ModelI, how TaskViewModel) string {
	displayingChildren := m.GetModel().Children
	if how.DisplayOnlyChildrenInProgress {
		displayingChildren = funk.Filter(m.GetModel().Children, func(c task.ModelI) bool {
			return c.GetProgress().InProgress
		}).([]task.ModelI)
	}
	if m.GetProgress().InProgress {
		if m.GetProgress().Ready == false {
			return m.ViewThis()
		}
		return lipgloss.JoinHorizontal(lipgloss.Left, funk.Map(displayingChildren, func(c task.ModelI) string {
			return c.View()
		}).([]string)...)
	}
	return "Not in progress"
}

func DisplayOnlyDoneOrInProgress(children []task.ModelI) string {

	displayingChildren := funk.Filter(children, func(c task.ModelI) bool {
		return c.GetProgress().InProgress || ((c.GetProgress().Completed || c.GetProgress().Skipped || c.GetProgress().Failed) && !c.GetProgress().HideOnDone)
	}).([]task.ModelI)

	return lipgloss.JoinVertical(lipgloss.Left, funk.Map(displayingChildren, func(c task.ModelI) string {
		return c.View()
	}).([]string)...)

}
