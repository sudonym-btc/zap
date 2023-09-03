package githubManager

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	"github.com/gookit/slog"
)

var client = github.NewClient(nil)

type GithubManager struct{}

func (*GithubManager) Name() string {
	return "github"
}

func (*GithubManager) FetchUser(name string) (*github.User, *github.Response, error) {
	slog.Debug("Fetching user", name)
	user, response, err := client.Users.Get(context.Background(), name)
	if err != nil {
		slog.Warn("Failed fetching user", err)
	} else {
		slog.Debug("Fetched user", user)
	}
	return user, response, err
}

func (*GithubManager) FetchRepo(name string) (*github.Repository, *github.Response, error) {
	slog.Debug("Fetching repo", name)
	return client.Repositories.Get(context.Background(), strings.Split(name, "/")[0], strings.Split(name, "/")[1])
}

func (*GithubManager) FetchOrg(name string) (*github.Organization, *github.Response, error) {
	slog.Debug("Fetching org", name)

	org, response, err := client.Organizations.Get(context.Background(), name)
	if err != nil {
		slog.Warn("Failed fetching org", err)
	} else {
		slog.Debug("Fetched org", org)
	}

	return org, response, err
}

func (*GithubManager) FetchOrgMaintainers(org github.Organization) ([]*github.User, *github.Response, error) {
	slog.Debug("Fetching org maintainers", org)
	req, err := client.NewRequest("GET", strings.Replace(org.GetPublicMembersURL(), "{/member}", "", -1), nil)
	if err != nil {
		slog.Warn("Failed fetching org maintainers", err)
		return nil, nil, err
	}
	ctx := context.Background()
	response := new([]*github.User)
	resp, err := client.Do(ctx, req, response)
	if err != nil {
		slog.Warn("Failed fetching org maintainers", err)
		return nil, nil, err
	}
	slog.Debug("Fetched org maintainers", *response)
	return *response, resp, nil
}

func (*GithubManager) FetchRepoMaintainers(repo github.Repository) ([]*github.User, *github.Response, error) {
	slog.Debug("Fetching repo maintainers", repo)
	req, err := client.NewRequest("GET", repo.GetContributorsURL(), nil)
	if err != nil {
		slog.Warn("Failed fetching repo maintainers", err)
		return nil, nil, err
	}
	ctx := context.Background()
	response := new([]*github.User)
	resp, err := client.Do(ctx, req, response)
	if err != nil {
		slog.Warn("Failed fetching repo maintainers", err)
		return nil, nil, err
	}
	slog.Debug("Fetched repo maintainers", *response)
	return *response, resp, nil
}
