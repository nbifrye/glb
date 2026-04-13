package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type getProjectArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
}

func registerProjectTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_project",
		Description: "Get detailed information about a GitLab project including description, visibility, default branch, and statistics.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getProjectArgs) (*mcp.CallToolResult, any, error) {
		p, _, err := client.Projects.GetProject(args.Project, nil)
		if err != nil {
			return textResult(fmt.Sprintf("Error getting project: %v", err)), nil, nil
		}

		data, _ := json.Marshal(p)
		return textResult(string(data)), nil, nil
	})
}
