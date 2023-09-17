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
)

type GithubRepoModel struct {
	task.Model
	tipModel TipModel
	repo     *github.Repository
	name     *string
}

func initialGithubRepoModel(name *string, tipModel TipModel) *GithubRepoModel {
	return &GithubRepoModel{
		tipModel: tipModel,
		name:     name,
		Model:    task.New(task.Model{}),
	}
}

func (m GithubRepoModel) Init() tea.Cmd {
	return nil
}

func (m GithubRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case repoMsg:
		if msg.id != m.GetId() {
			return m, nil
		}
		m.repo = &msg.content
		return m, fetchRepoUsers(m)
	case usersMsg:
		if msg.id != m.GetId() {
			return m, nil
		}
		for _, user := range msg.content {
			m.Children = append(m.Children, initialGithubUserModel(user.GetLogin(), false, m.tipModel))
		}
		cmds = append(cmds, task.Ready(m.GetId()))

	}

	return m, tea.Batch(cmds...)
}

func (m GithubRepoModel) View() string {

	return lipgloss.JoinVertical(lipgloss.Left, view.TitleStyle.Render("Github Repository"), indent.String(func() string {
		return lipgloss.JoinVertical(lipgloss.Left, view.ItemStyle(*m.name, ""), func() string {
			if m.Progress.InProgress && m.Progress.Ready == false {
				return view.PadVertical.Render(lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Loading repository..."))
			} else if len(m.Children) > 0 {
				return taskView.DisplayOnlyDoneOrInProgress(m.Children)
			} else if m.Progress.Failed {
				return lipgloss.JoinHorizontal(lipgloss.Left, view.Cross.Render(), view.ErrorMessageStyle("Failed fetching repo: "+fmt.Sprint(*m.Progress.Error)))
			}

			return "No users found"
		}())

	}(), view.DefaultIndent))

}

func (m GithubRepoModel) Job() tea.Cmd {

	return func() tea.Msg {

		repo, _, err := (&githubManager.GithubManager{}).FetchRepo(*m.name)
		if err != nil {
			return task.Failed(m.GetId(), &err)()
		}

		return repoMsg{id: m.GetId(), content: *repo}

	}
}

func fetchRepoUsers(m GithubRepoModel) tea.Cmd {
	return func() tea.Msg {

		users, _, err := (&githubManager.GithubManager{}).FetchRepoMaintainers(*m.repo)
		if err != nil {
			return task.Failed(m.GetId(), &err)()
		}

		return usersMsg{id: m.GetId(), content: users}
	}
}

type repoMsg struct {
	content github.Repository
	id      int
}
