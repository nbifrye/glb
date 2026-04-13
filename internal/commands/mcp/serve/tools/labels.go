package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/gitlabop"
)

type listLabelsArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	Search  string `json:"search,omitempty" jsonschema:"description=Search labels by keyword"`
	PerPage int64  `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 20\\, max: 100)"`
	Page    int64  `json:"page,omitempty" jsonschema:"description=Page number for pagination (default: 1)"`
}

func registerLabelTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_labels",
		Description: "List labels in a GitLab project. Returns label name, color, description, and issue/MR counts.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listLabelsArgs) (*mcp.CallToolResult, any, error) {
		labels, err := gitlabop.ListLabels(client, gitlabop.ListLabelsOptions{
			Project: args.Project,
			Search:  args.Search,
			PerPage: clampPerPage(args.PerPage, 20),
			Page:    args.Page,
		})
		if err != nil {
			return errorResult(fmt.Sprintf("Error listing labels: %v", err)), nil, nil
		}

		data, err := json.Marshal(labels)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})
}
