package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/gitlabop"
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
		p, err := gitlabop.GetProject(client, args.Project)
		if err != nil {
			return errorResult(fmt.Sprintf("Error getting project: %v", err)), nil, nil
		}

		data, err := json.Marshal(p)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})
}
