package tipView

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/github"
	"github.com/muesli/reflow/indent"
	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	taskView "github.com/sudonym-btc/zap/cmd/view/task"
	githubManager "github.com/sudonym-btc/zap/service/github-manager"
	"github.com/sudonym-btc/zap/service/tip"
)

type GithubUserModel struct {
	task.Model
	tipModel   TipModel
	showHeader bool
	user       *github.User
	name       string
}

func initialGithubUserModel(name string, showHeader bool, tipModel TipModel) *GithubUserModel {
	return &GithubUserModel{
		tipModel:   tipModel,
		showHeader: showHeader,
		name:       name,
		Model:      task.New(task.Model{}),
	}
}

func (m GithubUserModel) Init() tea.Cmd {
	return nil
}

func (m GithubUserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case userMsg:
		if msg.id != m.GetId() {
			return m, nil
		}
		m.user = &msg.content
		var maintainers []string

		if m.user.Email != nil {
			maintainers = append(maintainers, tip.ExtractEmails(*m.user.Email)...)
		}
		if m.user.Bio != nil {
			maintainers = append(maintainers, tip.ExtractEmails(*m.user.Bio)...)
		}

		for _, maintainer := range maintainers {
			m.Children = append(m.Children, initialMaintainerModel(maintainer, false, m.tipModel))
			cmds = append(cmds, m.Children[len(m.Children)-1].Init())
		}
		cmds = append(cmds, task.Ready(m.GetId()))

	}

	return m, tea.Batch(cmds...)
}

func (m GithubUserModel) View() string {

	content := func() string {
		return lipgloss.JoinVertical(lipgloss.Left, view.ItemStyle(m.name, ""), func() string {
			if m.Progress.InProgress {
				return view.PadVertical.Render(lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Loading user..."))
			} else if m.Progress.Failed {
				return lipgloss.JoinHorizontal(lipgloss.Left, view.Cross.Render(), view.ErrorMessageStyle("Failed fetching user: "+fmt.Sprint(*m.Progress.Error)))
			} else if len(m.Children) == 0 {
				return "No tippable addresses found"
			}
			return taskView.DisplayOnlyDoneOrInProgress(m.Children)
		}())
	}
	if m.showHeader {
		return lipgloss.JoinVertical(lipgloss.Left, view.TitleStyle.Render("Github User"), indent.String(content(), view.DefaultIndent))
	}
	return content()
}

func (m GithubUserModel) Job() tea.Cmd {

	return func() tea.Msg {

		user, _, err := (&githubManager.GithubManager{}).FetchUser(m.name)
		if err != nil {
			return task.Failed(m.GetId(), &err)()
		}
		return userMsg{id: m.GetId(), content: *user}

	}
}

type userMsg struct {
	content github.User
	id      int
}
