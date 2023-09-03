package tipView

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	taskView "github.com/sudonym-btc/zap/cmd/view/task"
	packageAnalyzer "github.com/sudonym-btc/zap/service/package-analyzer"
)

type PackageManagerModel struct {
	task.Model
	tipModel       TipModel
	packageManager packageAnalyzer.PackageManager
	packageName    *string
	packageVersion *string
}

func InitialPackageManagerModel(packageManager packageAnalyzer.PackageManager, packageName *string, packageVersion *string, tipModel TipModel) *PackageManagerModel {
	return &PackageManagerModel{
		tipModel:       tipModel,
		packageName:    packageName,
		packageVersion: packageVersion,
		packageManager: packageManager,
		Model:          task.New(task.Model{}),
	}
}

func (m PackageManagerModel) Init() tea.Cmd {
	return nil
}

func (m PackageManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case packagesMsg:
		if msg.id != m.GetId() {
			return m, nil
		}
		for _, pkg := range msg.content {
			m.Children = append(m.Children, initialPackageModel(m.packageManager, pkg, m.tipModel))
			cmds = append(cmds, m.Children[len(m.Children)-1].Init())
		}
		cmds = append(cmds, task.Ready(m.GetId()))
	}

	return m, tea.Batch(cmds...)
}

func (m PackageManagerModel) View() string {

	return lipgloss.JoinVertical(lipgloss.Left, view.TitleStyle.Render(m.packageManager.Name()), indent.String(func() string {
		if m.Progress.InProgress && m.Progress.Ready == false {
			return view.PadVertical.Render(lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Loading packages..."))
		} else if len(m.Children) > 0 {
			return taskView.DisplayOnlyDoneOrInProgress(m.Children)
		}

		return "No packages found"
	}(), view.DefaultIndent))

}

func (m PackageManagerModel) Job() tea.Cmd {

	return func() tea.Msg {

		if m.packageName != nil && *m.packageName != "" {
			packages := []*packageAnalyzer.PackageInfo{{
				Name:    *m.packageName,
				Version: *m.packageVersion,
			}}
			return packagesMsg{id: m.GetId(), content: packages}

		} else {
			packages, err := m.packageManager.FetchPackages()
			if err != nil {
				return task.Failed(m.GetId(), &err)()
			}
			return packagesMsg{id: m.GetId(), content: packages}

		}
	}
}

type packagesMsg struct {
	content []*packageAnalyzer.PackageInfo
	id      int
}
