package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type listPipelineArtifactsArgs struct {
	Project    string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	PipelineID int64 `json:"pipeline_id" jsonschema:"required,description=Pipeline ID"`
	PerPage    int64 `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 50)"`
}

type jobArtifactInfo struct {
	JobID     int64                `json:"job_id"`
	JobName   string               `json:"job_name"`
	Status    string               `json:"status"`
	Stage     string               `json:"stage"`
	Artifacts []gitlab.JobArtifact `json:"artifacts"`
}

func registerArtifactTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pipeline_artifacts",
		Description: "List jobs and their artifacts for a pipeline. Returns job names, statuses, and artifact information. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listPipelineArtifactsArgs) (*mcp.CallToolResult, any, error) {
		perPage := args.PerPage
		if perPage == 0 {
			perPage = 50
		}

		jobs, _, err := client.Jobs.ListPipelineJobs(args.Project, args.PipelineID, &gitlab.ListJobsOptions{
			ListOptions: gitlab.ListOptions{PerPage: perPage},
		})
		if err != nil {
			return textResult(fmt.Sprintf("Error listing jobs: %v", err)), nil, nil
		}

		results := make([]jobArtifactInfo, 0)
		for _, job := range jobs {
			if len(job.Artifacts) > 0 {
				results = append(results, jobArtifactInfo{
					JobID:     job.ID,
					JobName:   job.Name,
					Status:    job.Status,
					Stage:     job.Stage,
					Artifacts: job.Artifacts,
				})
			}
		}

		data, _ := json.Marshal(results)
		return textResult(string(data)), nil, nil
	})
}
