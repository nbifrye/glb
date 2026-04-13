package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type listPipelinesArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	Status  string `json:"status,omitempty" jsonschema:"description=Filter by status: running\\, pending\\, success\\, failed\\, canceled"`
	Ref     string `json:"ref,omitempty" jsonschema:"description=Filter by git ref (branch or tag name)"`
	PerPage int64  `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 20)"`
}

type getPipelineArgs struct {
	Project    string `json:"project" jsonschema:"required,description=Project path"`
	PipelineID int64  `json:"pipeline_id" jsonschema:"required,description=Pipeline ID"`
}

func registerPipelineTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pipelines",
		Description: "List CI/CD pipelines for a project. Returns pipeline ID, status, ref, and creation time.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listPipelinesArgs) (*mcp.CallToolResult, any, error) {
		perPage := args.PerPage
		if perPage == 0 {
			perPage = 20
		}
		opts := &gitlab.ListProjectPipelinesOptions{
			ListOptions: gitlab.ListOptions{PerPage: perPage},
		}
		if args.Status != "" {
			s := gitlab.BuildStateValue(args.Status)
			opts.Status = &s
		}
		if args.Ref != "" {
			opts.Ref = gitlab.Ptr(args.Ref)
		}

		pipelines, _, err := client.Pipelines.ListProjectPipelines(args.Project, opts)
		if err != nil {
			return textResult(fmt.Sprintf("Error listing pipelines: %v", err)), nil, nil
		}

		data, _ := json.Marshal(pipelines)
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pipeline",
		Description: "Get detailed information about a specific pipeline including status, duration, and coverage.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getPipelineArgs) (*mcp.CallToolResult, any, error) {
		pipeline, _, err := client.Pipelines.GetPipeline(args.Project, args.PipelineID)
		if err != nil {
			return textResult(fmt.Sprintf("Error getting pipeline: %v", err)), nil, nil
		}

		data, _ := json.Marshal(pipeline)
		return textResult(string(data)), nil, nil
	})
}
