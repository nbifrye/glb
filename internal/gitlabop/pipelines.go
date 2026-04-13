package gitlabop

import (
	"bytes"
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type ListPipelinesOptions struct {
	Project string
	Status  string
	Ref     string
	PerPage int64
	Page    int64
}

func ListPipelines(client *gitlab.Client, opts ListPipelinesOptions) ([]*gitlab.PipelineInfo, error) {
	perPage := opts.PerPage
	if perPage <= 0 {
		perPage = int64(DefaultPerPage)
	}
	apiOpts := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{PerPage: perPage},
	}
	if opts.Page > 0 {
		apiOpts.ListOptions.Page = opts.Page
	}
	if opts.Status != "" {
		s := gitlab.BuildStateValue(opts.Status)
		apiOpts.Status = &s
	}
	if opts.Ref != "" {
		apiOpts.Ref = gitlab.Ptr(opts.Ref)
	}

	pipelines, _, err := client.Pipelines.ListProjectPipelines(opts.Project, apiOpts)
	if err != nil {
		return nil, fmt.Errorf("listing pipelines: %w", err)
	}
	return pipelines, nil
}

func GetPipeline(client *gitlab.Client, project string, id int64) (*gitlab.Pipeline, error) {
	pipeline, _, err := client.Pipelines.GetPipeline(project, id)
	if err != nil {
		return nil, fmt.Errorf("getting pipeline: %w", err)
	}
	return pipeline, nil
}

type ListPipelineJobsOptions struct {
	Project    string
	PipelineID int64
	Scope      string
	PerPage    int64
	Page       int64
}

func ListPipelineJobs(client *gitlab.Client, opts ListPipelineJobsOptions) ([]*gitlab.Job, error) {
	perPage := opts.PerPage
	if perPage <= 0 {
		perPage = int64(DefaultPerPage)
	}
	apiOpts := &gitlab.ListJobsOptions{
		ListOptions: gitlab.ListOptions{PerPage: perPage},
	}
	if opts.Page > 0 {
		apiOpts.ListOptions.Page = opts.Page
	}
	if opts.Scope != "" {
		apiOpts.Scope = &[]gitlab.BuildStateValue{gitlab.BuildStateValue(opts.Scope)}
	}

	jobs, _, err := client.Jobs.ListPipelineJobs(opts.Project, opts.PipelineID, apiOpts)
	if err != nil {
		return nil, fmt.Errorf("listing pipeline jobs: %w", err)
	}
	return jobs, nil
}

func GetJobLog(client *gitlab.Client, project string, jobID int64) (string, error) {
	reader, _, err := client.Jobs.GetTraceFile(project, jobID)
	if err != nil {
		return "", fmt.Errorf("getting job log: %w", err)
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return "", fmt.Errorf("reading job log: %w", err)
	}
	return buf.String(), nil
}
