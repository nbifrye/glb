package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type compareRefsArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	From    string `json:"from" jsonschema:"required,description=Source ref (branch name\\, tag\\, or commit SHA)"`
	To      string `json:"to" jsonschema:"required,description=Target ref (branch name\\, tag\\, or commit SHA)"`
}

func registerCompareTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "compare_refs",
		Description: "Compare two branches, tags, or commits in a GitLab repository. Returns the list of commits and file diffs between the two refs. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args compareRefsArgs) (*mcp.CallToolResult, any, error) {
		compare, _, err := client.Repositories.Compare(args.Project, &gitlab.CompareOptions{
			From: gitlab.Ptr(args.From),
			To:   gitlab.Ptr(args.To),
		})
		if err != nil {
			return errorResult(fmt.Sprintf("Error comparing refs: %v", err)), nil, nil
		}

		data, err := json.Marshal(compare)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})
}
