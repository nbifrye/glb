package gitlabop

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type ListMergeRequestsOptions struct {
	Project string
	State   string
	Labels  []string
	PerPage int64
	Page    int64
}

func ListMergeRequests(client *gitlab.Client, opts ListMergeRequestsOptions) ([]*gitlab.BasicMergeRequest, error) {
	perPage := opts.PerPage
	if perPage <= 0 {
		perPage = int64(DefaultPerPage)
	}
	apiOpts := &gitlab.ListProjectMergeRequestsOptions{
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

	mrs, _, err := client.MergeRequests.ListProjectMergeRequests(opts.Project, apiOpts)
	if err != nil {
		return nil, fmt.Errorf("listing merge requests: %w", err)
	}
	return mrs, nil
}

func GetMergeRequest(client *gitlab.Client, project string, iid int64) (*gitlab.MergeRequest, error) {
	mr, _, err := client.MergeRequests.GetMergeRequest(project, iid, nil)
	if err != nil {
		return nil, fmt.Errorf("getting merge request: %w", err)
	}
	return mr, nil
}

type CreateMergeRequestOptions struct {
	Project      string
	Title        string
	SourceBranch string
	TargetBranch string
	Description  string
	Labels       []string
	Draft        bool
}

func CreateMergeRequest(client *gitlab.Client, opts CreateMergeRequestOptions) (*gitlab.MergeRequest, error) {
	title := opts.Title
	if opts.Draft {
		title = "Draft: " + title
	}

	apiOpts := &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.Ptr(title),
		SourceBranch: gitlab.Ptr(opts.SourceBranch),
		TargetBranch: gitlab.Ptr(opts.TargetBranch),
	}
	if opts.Description != "" {
		apiOpts.Description = gitlab.Ptr(opts.Description)
	}
	if len(opts.Labels) > 0 {
		apiOpts.Labels = (*gitlab.LabelOptions)(&opts.Labels)
	}

	mr, _, err := client.MergeRequests.CreateMergeRequest(opts.Project, apiOpts)
	if err != nil {
		return nil, fmt.Errorf("creating merge request: %w", err)
	}
	return mr, nil
}

type GetMergeRequestDiffResult struct {
	*gitlab.MergeRequestDiffVersion
}

func GetMergeRequestDiff(client *gitlab.Client, project string, iid int64) (*gitlab.MergeRequestDiffVersion, error) {
	versions, _, err := client.MergeRequests.GetMergeRequestDiffVersions(project, iid, &gitlab.GetMergeRequestDiffVersionsOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting diff versions: %w", err)
	}
	if len(versions) == 0 {
		return nil, nil
	}

	version, _, err := client.MergeRequests.GetSingleMergeRequestDiffVersion(project, iid, versions[0].ID, &gitlab.GetSingleMergeRequestDiffVersionOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting diff: %w", err)
	}
	return version, nil
}

type MergeMergeRequestOptions struct {
	Project      string
	IID          int64
	Squash       bool
	RemoveBranch bool
}

func MergeMergeRequest(client *gitlab.Client, opts MergeMergeRequestOptions) (*gitlab.MergeRequest, error) {
	apiOpts := &gitlab.AcceptMergeRequestOptions{}
	if opts.Squash {
		apiOpts.Squash = gitlab.Ptr(true)
	}
	if opts.RemoveBranch {
		apiOpts.ShouldRemoveSourceBranch = gitlab.Ptr(true)
	}

	mr, _, err := client.MergeRequests.AcceptMergeRequest(opts.Project, opts.IID, apiOpts)
	if err != nil {
		return nil, fmt.Errorf("merging MR: %w", err)
	}
	return mr, nil
}

func ApproveMergeRequest(client *gitlab.Client, project string, iid int64) (*gitlab.MergeRequestApprovals, error) {
	approval, _, err := client.MergeRequestApprovals.ApproveMergeRequest(project, iid, nil)
	if err != nil {
		return nil, fmt.Errorf("approving merge request: %w", err)
	}
	return approval, nil
}

func UnapproveMergeRequest(client *gitlab.Client, project string, iid int64) error {
	_, err := client.MergeRequestApprovals.UnapproveMergeRequest(project, iid)
	if err != nil {
		return fmt.Errorf("unapproving merge request: %w", err)
	}
	return nil
}

func AddMergeRequestNote(client *gitlab.Client, project string, iid int64, body string) (*gitlab.Note, error) {
	note, _, err := client.Notes.CreateMergeRequestNote(project, iid, &gitlab.CreateMergeRequestNoteOptions{
		Body: gitlab.Ptr(body),
	})
	if err != nil {
		return nil, fmt.Errorf("adding note: %w", err)
	}
	return note, nil
}
