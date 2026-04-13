package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type batchUpdateArgs struct {
	Project      string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	ResourceType string `json:"resource_type" jsonschema:"required,description=Resource type: 'issue' or 'mr'"`
	ResourceIDs  string `json:"resource_ids" jsonschema:"required,description=Comma-separated list of resource IIDs to update"`
	AddLabels    string `json:"add_labels,omitempty" jsonschema:"description=Comma-separated labels to add"`
	RemoveLabels string `json:"remove_labels,omitempty" jsonschema:"description=Comma-separated labels to remove"`
	StateEvent   string `json:"state_event,omitempty" jsonschema:"description=State change: 'close' or 'reopen'"`
	MilestoneID  int64  `json:"milestone_id,omitempty" jsonschema:"description=Milestone ID to set"`
}

type updateResult struct {
	ID      int64  `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func registerBatchTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "batch_update",
		Description: "Bulk update multiple issues or merge requests at once. Supports changing labels, state, and milestone. Particularly useful for AI agents triaging multiple items. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(true)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args batchUpdateArgs) (*mcp.CallToolResult, any, error) {
		ids := parseIntList(args.ResourceIDs)
		if len(ids) == 0 {
			return textResult("Error: no valid resource IDs provided"), nil, nil
		}

		results := make([]updateResult, 0, len(ids))

		for _, id := range ids {
			var err error
			switch args.ResourceType {
			case "issue":
				opts := &gitlab.UpdateIssueOptions{}
				if args.AddLabels != "" {
					labels := splitLabels(args.AddLabels)
					opts.AddLabels = (*gitlab.LabelOptions)(&labels)
				}
				if args.RemoveLabels != "" {
					labels := splitLabels(args.RemoveLabels)
					opts.RemoveLabels = (*gitlab.LabelOptions)(&labels)
				}
				if args.StateEvent != "" {
					opts.StateEvent = gitlab.Ptr(args.StateEvent)
				}
				if args.MilestoneID > 0 {
					opts.MilestoneID = gitlab.Ptr(args.MilestoneID)
				}
				_, _, err = client.Issues.UpdateIssue(args.Project, id, opts)

			case "mr":
				opts := &gitlab.UpdateMergeRequestOptions{}
				if args.AddLabels != "" {
					labels := splitLabels(args.AddLabels)
					opts.AddLabels = (*gitlab.LabelOptions)(&labels)
				}
				if args.RemoveLabels != "" {
					labels := splitLabels(args.RemoveLabels)
					opts.RemoveLabels = (*gitlab.LabelOptions)(&labels)
				}
				if args.StateEvent != "" {
					opts.StateEvent = gitlab.Ptr(args.StateEvent)
				}
				if args.MilestoneID > 0 {
					opts.MilestoneID = gitlab.Ptr(args.MilestoneID)
				}
				_, _, err = client.MergeRequests.UpdateMergeRequest(args.Project, id, opts)

			default:
				return textResult("Error: resource_type must be 'issue' or 'mr'"), nil, nil
			}

			if err != nil {
				results = append(results, updateResult{ID: id, Status: "error", Message: err.Error()})
			} else {
				results = append(results, updateResult{ID: id, Status: "updated"})
			}
		}

		data, _ := json.Marshal(results)
		return textResult(string(data)), nil, nil
	})
}

func parseIntList(s string) []int64 {
	parts := strings.Split(s, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		var n int64
		if _, err := fmt.Sscanf(p, "%d", &n); err == nil {
			ids = append(ids, n)
		}
	}
	return ids
}
