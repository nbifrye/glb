package gitlabop

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type ListIssuesOptions struct {
	Project  string
	State    string
	Labels   []string
	Assignee string
	Search   string
	PerPage  int64
	Page     int64
}

func ListIssues(client *gitlab.Client, opts ListIssuesOptions) ([]*gitlab.Issue, error) {
	perPage := opts.PerPage
	if perPage <= 0 {
		perPage = int64(DefaultPerPage)
	}
	apiOpts := &gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{PerPage: perPage},
	}
	if opts.Page > 0 {
		apiOpts.ListOptions.Page = opts.Page
	}
	if opts.State != "" {
		apiOpts.State = gitlab.Ptr(opts.State)
	}
	if len(opts.Labels) > 0 {
		apiOpts.Labels = (*gitlab.LabelOptions)(&opts.Labels)
	}
	if opts.Search != "" {
		apiOpts.Search = gitlab.Ptr(opts.Search)
	}
	if opts.Assignee != "" {
		users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{Username: gitlab.Ptr(opts.Assignee)})
		if err != nil {
			return nil, fmt.Errorf("resolving assignee %q: %w", opts.Assignee, err)
		}
		if len(users) == 0 {
			return nil, fmt.Errorf("user %q not found", opts.Assignee)
		}
		apiOpts.AssigneeID = gitlab.AssigneeID(users[0].ID)
	}

	issues, _, err := client.Issues.ListProjectIssues(opts.Project, apiOpts)
	if err != nil {
		return nil, fmt.Errorf("listing issues: %w", err)
	}
	return issues, nil
}

func GetIssue(client *gitlab.Client, project string, iid int64) (*gitlab.Issue, error) {
	issue, _, err := client.Issues.GetIssue(project, iid)
	if err != nil {
		return nil, fmt.Errorf("getting issue: %w", err)
	}
	return issue, nil
}

type CreateIssueOptions struct {
	Project     string
	Title       string
	Description string
	Labels      []string
	AssigneeIDs []int64
}

func CreateIssue(client *gitlab.Client, opts CreateIssueOptions) (*gitlab.Issue, error) {
	apiOpts := &gitlab.CreateIssueOptions{
		Title: gitlab.Ptr(opts.Title),
	}
	if opts.Description != "" {
		apiOpts.Description = gitlab.Ptr(opts.Description)
	}
	if len(opts.Labels) > 0 {
		apiOpts.Labels = (*gitlab.LabelOptions)(&opts.Labels)
	}
	if len(opts.AssigneeIDs) > 0 {
		apiOpts.AssigneeIDs = &opts.AssigneeIDs
	}

	issue, _, err := client.Issues.CreateIssue(opts.Project, apiOpts)
	if err != nil {
		return nil, fmt.Errorf("creating issue: %w", err)
	}
	return issue, nil
}

func CloseIssue(client *gitlab.Client, project string, iid int64) (*gitlab.Issue, error) {
	issue, _, err := client.Issues.UpdateIssue(project, iid, &gitlab.UpdateIssueOptions{
		StateEvent: gitlab.Ptr("close"),
	})
	if err != nil {
		return nil, fmt.Errorf("closing issue: %w", err)
	}
	return issue, nil
}

func AddIssueNote(client *gitlab.Client, project string, iid int64, body string) (*gitlab.Note, error) {
	note, _, err := client.Notes.CreateIssueNote(project, iid, &gitlab.CreateIssueNoteOptions{
		Body: gitlab.Ptr(body),
	})
	if err != nil {
		return nil, fmt.Errorf("adding note: %w", err)
	}
	return note, nil
}
