package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type getMRConflictsArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	MRID    int64  `json:"mr_id" jsonschema:"required,description=Merge request IID"`
}

type conflictInfo struct {
	MRID                int64  `json:"mr_iid"`
	Title               string `json:"title"`
	HasConflicts        bool   `json:"has_conflicts"`
	SourceBranch        string `json:"source_branch"`
	TargetBranch        string `json:"target_branch"`
	DetailedMergeStatus string `json:"detailed_merge_status"`
}

func registerMRConflictTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_mr_conflicts",
		Description: "Get merge request conflict information including whether conflicts exist and merge status. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getMRConflictsArgs) (*mcp.CallToolResult, any, error) {
		mr, _, err := client.MergeRequests.GetMergeRequest(args.Project, args.MRID, nil)
		if err != nil {
			return textResult(fmt.Sprintf("Error getting MR: %v", err)), nil, nil
		}

		info := conflictInfo{
			MRID:                mr.IID,
			Title:               mr.Title,
			HasConflicts:        mr.HasConflicts,
			SourceBranch:        mr.SourceBranch,
			TargetBranch:        mr.TargetBranch,
			DetailedMergeStatus: mr.DetailedMergeStatus,
		}

		data, _ := json.Marshal(info)
		return textResult(string(data)), nil, nil
	})
}
