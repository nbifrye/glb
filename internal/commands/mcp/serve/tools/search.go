package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type searchCodeArgs struct {
	Query   string `json:"query" jsonschema:"required,description=Search query string"`
	Project string `json:"project,omitempty" jsonschema:"description=Limit search to a specific project path"`
	Group   string `json:"group,omitempty" jsonschema:"description=Limit search to a specific group path"`
	PerPage int64  `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 20)"`
}

func registerSearchTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_code",
		Description: "Search code across GitLab projects or within a specific project/group. Returns matching file paths and code snippets. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args searchCodeArgs) (*mcp.CallToolResult, any, error) {
		perPage := args.PerPage
		if perPage == 0 {
			perPage = 20
		}
		opts := &gitlab.SearchOptions{
			ListOptions: gitlab.ListOptions{PerPage: perPage},
		}

		if args.Project != "" {
			results, _, err := client.Search.BlobsByProject(args.Project, args.Query, opts)
			if err != nil {
				return errorResult(fmt.Sprintf("Error searching: %v", err)), nil, nil
			}
			data, err := json.Marshal(results)
			if err != nil {
				return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
			}
			return textResult(string(data)), nil, nil
		}

		if args.Group != "" {
			results, _, err := client.Search.BlobsByGroup(args.Group, args.Query, opts)
			if err != nil {
				return errorResult(fmt.Sprintf("Error searching: %v", err)), nil, nil
			}
			data, err := json.Marshal(results)
			if err != nil {
				return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
			}
			return textResult(string(data)), nil, nil
		}

		results, _, err := client.Search.Blobs(args.Query, opts)
		if err != nil {
			return errorResult(fmt.Sprintf("Error searching: %v", err)), nil, nil
		}
		data, err := json.Marshal(results)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})
}
