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

type GithubOrganizationModel struct {
	task.Model
	tipModel     TipModel
	organization *github.Organization
	name         *string
}

func initialGithubOrganizationModel(name *string, tipModel TipModel) *GithubOrganizationModel {
	return &GithubOrganizationModel{
		tipModel: tipModel,
		name:     name,
		Model:    task.New(task.Model{}),
	}
}

func (m GithubOrganizationModel) Init() tea.Cmd {
	return nil
}

func (m GithubOrganizationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case organizationMsg:
		if msg.id != m.GetId() {
			return m, nil
		}
		m.organization = &msg.content
		return m, fetchOrganizationUsers(m)
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

func (m GithubOrganizationModel) View() string {

	return lipgloss.JoinVertical(lipgloss.Left, view.TitleStyle.Render("Github Organization"), indent.String(func() string {
		return lipgloss.JoinVertical(lipgloss.Left, view.ItemStyle(*m.name, ""), func() string {
			if m.Progress.InProgress && m.Progress.Ready == false {
				return view.PadVertical.Render(lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Loading organization..."))
			} else if len(m.Children) > 0 {
				return taskView.DisplayOnlyDoneOrInProgress(m.Children)
			} else if m.Progress.Failed {
				return lipgloss.JoinHorizontal(lipgloss.Left, view.Cross.Render(), view.ErrorMessageStyle("Failed fetching organisation: "+fmt.Sprint(*m.Progress.Error)))
			}

			return "No public users found"
		}())

	}(), view.DefaultIndent))

}

func (m GithubOrganizationModel) Job() tea.Cmd {

	return func() tea.Msg {

		org, _, err := (&githubManager.GithubManager{}).FetchOrg(*m.name)
		if err != nil {
			return task.Failed(m.GetId(), &err)()
		}

		return organizationMsg{id: m.GetId(), content: *org}

	}
}

func fetchOrganizationUsers(m GithubOrganizationModel) tea.Cmd {
	return func() tea.Msg {

		users, _, err := (&githubManager.GithubManager{}).FetchOrgMaintainers(*m.organization)
		if err != nil {
			return task.Failed(m.GetId(), &err)()
		}

		return usersMsg{id: m.GetId(), content: users}
	}
}

type organizationMsg struct {
	content github.Organization
	id      int
}

type usersMsg struct {
	content []*github.User
	id      int
}
