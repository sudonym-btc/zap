package tipView

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	taskView "github.com/sudonym-btc/zap/cmd/view/task"
	packageAnalyzer "github.com/sudonym-btc/zap/service/package-analyzer"
	"github.com/sudonym-btc/zap/service/tip"
)

type packageModel struct {
	task.Model
	tipModel       TipModel
	packageManager packageAnalyzer.PackageManager
	packageInfo    packageAnalyzer.PackageInfo
}

func initialPackageModel(packageManager packageAnalyzer.PackageManager, info *packageAnalyzer.PackageInfo, tipModel TipModel) *packageModel {
	return &packageModel{
		tipModel:       tipModel,
		packageInfo:    *info,
		packageManager: packageManager,
		Model:          task.New(task.Model{}),
	}
}

func (m packageModel) Init() tea.Cmd {
	return nil
}

func (m packageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case packageInfoMsg:
		if msg.id != m.GetId() {
			return m, nil
		}
		m.packageInfo = *msg.content
		if len(m.packageInfo.Funding) > 0 {
			for _, funding := range m.packageInfo.Funding {
				if len(tip.ExtractEmails(funding.URL)) > 0 {
					m.Children = append(m.Children, initialMaintainerModel(tip.ExtractEmails(funding.URL)[0], false, m.tipModel))
					cmds = append(cmds, m.Children[len(m.Children)-1].Init())
				}

			}
		}
		for _, maintainer := range m.packageInfo.Maintainers {
			m.Children = append(m.Children, initialMaintainerModel(maintainer, false, m.tipModel))
			cmds = append(cmds, m.Children[len(m.Children)-1].Init())
		}
		cmds = append(cmds, task.Ready(m.GetId()))
	}

	return m, tea.Batch(cmds...)
}

func (m packageModel) View() string {

	return lipgloss.JoinVertical(lipgloss.Left, view.ItemStyle(m.packageInfo.Name, m.packageInfo.Version), func() string {

		if m.Progress.InProgress && m.Progress.Ready == false {
			return lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Loading package info...")
		} else if len(m.Children) > 0 {
			return taskView.DisplayOnlyDoneOrInProgress(m.Children)
		} else if m.Progress.Failed {
			return lipgloss.JoinHorizontal(lipgloss.Left, view.Cross.Render(), view.ErrorMessageStyle("Failed loading package: "+fmt.Sprint(*m.Progress.Error)))
		}
		return indent.String("Waiting to process...", view.DefaultIndent)
	}())

}

func (m packageModel) Job() tea.Cmd {

	return func() tea.Msg {

		info, err := m.packageManager.FetchInfo(m.packageInfo.Name, m.packageInfo.Version)
		if err != nil {
			return task.Failed(m.GetId(), &err)()
		}
		return packageInfoMsg{id: m.GetId(), content: info}
	}

}

type packageInfoMsg struct {
	id      int
	content *packageAnalyzer.PackageInfo
}
