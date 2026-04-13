package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/gitlabop"
)

type listPipelinesArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	Status  string `json:"status,omitempty" jsonschema:"description=Filter by status: running\\, pending\\, success\\, failed\\, canceled"`
	Ref     string `json:"ref,omitempty" jsonschema:"description=Filter by git ref (branch or tag name)"`
	PerPage int64  `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 20\\, max: 100)"`
	Page    int64  `json:"page,omitempty" jsonschema:"description=Page number for pagination (default: 1)"`
}

type getPipelineArgs struct {
	Project    string `json:"project" jsonschema:"required,description=Project path"`
	PipelineID int64  `json:"pipeline_id" jsonschema:"required,description=Pipeline ID"`
}

type listPipelineJobsArgs struct {
	Project    string `json:"project" jsonschema:"required,description=Project path"`
	PipelineID int64  `json:"pipeline_id" jsonschema:"required,description=Pipeline ID"`
	Scope      string `json:"scope,omitempty" jsonschema:"description=Filter by scope: created\\, pending\\, running\\, failed\\, success\\, canceled\\, skipped"`
	PerPage    int64  `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 20\\, max: 100)"`
	Page       int64  `json:"page,omitempty" jsonschema:"description=Page number for pagination (default: 1)"`
}

type getJobLogArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path"`
	JobID   int64  `json:"job_id" jsonschema:"required,description=Job ID"`
}

func registerPipelineTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pipelines",
		Description: "List CI/CD pipelines for a project. Returns pipeline ID, status, ref, and creation time.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listPipelinesArgs) (*mcp.CallToolResult, any, error) {
		pipelines, err := gitlabop.ListPipelines(client, gitlabop.ListPipelinesOptions{
			Project: args.Project,
			Status:  args.Status,
			Ref:     args.Ref,
			PerPage: clampPerPage(args.PerPage, 20),
			Page:    args.Page,
		})
		if err != nil {
			return errorResult(fmt.Sprintf("Error listing pipelines: %v", err)), nil, nil
		}

		data, err := json.Marshal(pipelines)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pipeline",
		Description: "Get detailed information about a specific pipeline including status, duration, and coverage.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getPipelineArgs) (*mcp.CallToolResult, any, error) {
		pipeline, err := gitlabop.GetPipeline(client, args.Project, args.PipelineID)
		if err != nil {
			return errorResult(fmt.Sprintf("Error getting pipeline: %v", err)), nil, nil
		}

		data, err := json.Marshal(pipeline)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pipeline_jobs",
		Description: "List jobs for a specific pipeline. Returns job ID, name, stage, and status.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listPipelineJobsArgs) (*mcp.CallToolResult, any, error) {
		jobs, err := gitlabop.ListPipelineJobs(client, gitlabop.ListPipelineJobsOptions{
			Project:    args.Project,
			PipelineID: args.PipelineID,
			Scope:      args.Scope,
			PerPage:    clampPerPage(args.PerPage, 20),
			Page:       args.Page,
		})
		if err != nil {
			return errorResult(fmt.Sprintf("Error listing pipeline jobs: %v", err)), nil, nil
		}

		data, err := json.Marshal(jobs)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_job_log",
		Description: "Get the log/trace output of a specific CI/CD job.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getJobLogArgs) (*mcp.CallToolResult, any, error) {
		logOutput, err := gitlabop.GetJobLog(client, args.Project, args.JobID)
		if err != nil {
			return errorResult(fmt.Sprintf("Error getting job log: %v", err)), nil, nil
		}

		return textResult(logOutput), nil, nil
	})
}
