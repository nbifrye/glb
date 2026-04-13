package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type addTimeSpentArgs struct {
	Project      string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	ResourceType string `json:"resource_type" jsonschema:"required,description=Resource type: 'issue' or 'mr'"`
	ResourceID   int64  `json:"resource_id" jsonschema:"required,description=Resource IID"`
	Duration     string `json:"duration" jsonschema:"required,description=Time duration (e.g. '1h30m'\\, '2d'\\, '30m')"`
}

type setTimeEstimateArgs struct {
	Project      string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	ResourceType string `json:"resource_type" jsonschema:"required,description=Resource type: 'issue' or 'mr'"`
	ResourceID   int64  `json:"resource_id" jsonschema:"required,description=Resource IID"`
	Duration     string `json:"duration" jsonschema:"required,description=Time estimate (e.g. '4h'\\, '1d'\\, '2h30m')"`
}

func registerTimetrackingTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_time_spent",
		Description: "Add time spent on an issue or merge request. Uses GitLab's time tracking format. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(false)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addTimeSpentArgs) (*mcp.CallToolResult, any, error) {
		switch args.ResourceType {
		case "issue":
			stats, _, err := client.Issues.AddSpentTime(args.Project, args.ResourceID, &gitlab.AddSpentTimeOptions{
				Duration: gitlab.Ptr(args.Duration),
			})
			if err != nil {
				return errorResult(fmt.Sprintf("Error adding time: %v", err)), nil, nil
			}
			data, err := json.Marshal(stats)
			if err != nil {
				return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
			}
			return textResult(string(data)), nil, nil

		case "mr":
			stats, _, err := client.MergeRequests.AddSpentTime(args.Project, args.ResourceID, &gitlab.AddSpentTimeOptions{
				Duration: gitlab.Ptr(args.Duration),
			})
			if err != nil {
				return errorResult(fmt.Sprintf("Error adding time: %v", err)), nil, nil
			}
			data, err := json.Marshal(stats)
			if err != nil {
				return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
			}
			return textResult(string(data)), nil, nil

		default:
			return errorResult("Error: resource_type must be 'issue' or 'mr'"), nil, nil
		}
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "set_time_estimate",
		Description: "Set a time estimate on an issue or merge request. Uses GitLab's time tracking format. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(false)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args setTimeEstimateArgs) (*mcp.CallToolResult, any, error) {
		switch args.ResourceType {
		case "issue":
			stats, _, err := client.Issues.SetTimeEstimate(args.Project, args.ResourceID, &gitlab.SetTimeEstimateOptions{
				Duration: gitlab.Ptr(args.Duration),
			})
			if err != nil {
				return errorResult(fmt.Sprintf("Error setting estimate: %v", err)), nil, nil
			}
			data, err := json.Marshal(stats)
			if err != nil {
				return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
			}
			return textResult(string(data)), nil, nil

		case "mr":
			stats, _, err := client.MergeRequests.SetTimeEstimate(args.Project, args.ResourceID, &gitlab.SetTimeEstimateOptions{
				Duration: gitlab.Ptr(args.Duration),
			})
			if err != nil {
				return errorResult(fmt.Sprintf("Error setting estimate: %v", err)), nil, nil
			}
			data, err := json.Marshal(stats)
			if err != nil {
				return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
			}
			return textResult(string(data)), nil, nil

		default:
			return errorResult("Error: resource_type must be 'issue' or 'mr'"), nil, nil
		}
	})
}
