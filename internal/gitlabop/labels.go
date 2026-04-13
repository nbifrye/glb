package gitlabop

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type ListLabelsOptions struct {
	Project string
	Search  string
	PerPage int64
	Page    int64
}

func ListLabels(client *gitlab.Client, opts ListLabelsOptions) ([]*gitlab.Label, error) {
	perPage := opts.PerPage
	if perPage <= 0 {
		perPage = int64(DefaultPerPage)
	}
	apiOpts := &gitlab.ListLabelsOptions{
		ListOptions: gitlab.ListOptions{PerPage: perPage},
		WithCounts:  gitlab.Ptr(true),
	}
	if opts.Page > 0 {
		apiOpts.ListOptions.Page = opts.Page
	}
	if opts.Search != "" {
		apiOpts.Search = gitlab.Ptr(opts.Search)
	}

	labels, _, err := client.Labels.ListLabels(opts.Project, apiOpts)
	if err != nil {
		return nil, fmt.Errorf("listing labels: %w", err)
	}
	return labels, nil
}
