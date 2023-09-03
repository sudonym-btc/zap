package tipView

import (
	"github.com/muesli/reflow/indent"
	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	taskView "github.com/sudonym-btc/zap/cmd/view/task"
	"github.com/sudonym-btc/zap/service/tip"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type maintainerModel struct {
	task.Model
	address    string
	showHeader bool
	tipModel   TipModel
}

type detailConfirmParams struct {
	text      *string
	textLimit *int
	amount    *int
}

func initialMaintainerModel(address string, showHeader bool, tipModel TipModel) *maintainerModel {
	m := &maintainerModel{
		showHeader: showHeader,
		tipModel:   tipModel,
		address:    tip.ExtractEmails(address)[0],
		Model:      task.New(task.Model{Progress: &task.TaskProgress{ShouldCompleteOnFirstSubtaskComplete: true}}),
	}
	m.Children = append(m.Children, initialMaintainerLnFlowModel(m))
	if tipModel.sendEmails {
		m.Children = append(m.Children, initialMaintainerEmailFlowModel(m))
	}
	return m
}

func (m maintainerModel) Job() tea.Cmd {
	return nil
}

func (m maintainerModel) Init() tea.Cmd {
	return nil
}

func (m maintainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	return m, tea.Batch(cmds...)
}

func (m maintainerModel) View() string {

	title := view.Faint.Render(m.address)

	content := func() string {
		return lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Italic(true).Render(title),
			taskView.DisplayOnlyDoneOrInProgress(m.Children))
	}

	if m.showHeader {
		return lipgloss.JoinVertical(lipgloss.Left, view.TitleStyle.Render("Tip"), indent.String(content(), view.DefaultIndent))
	} else {
		return content()
	}

}
