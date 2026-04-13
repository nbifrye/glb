package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type getRepoFileArgs struct {
	Project  string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	FilePath string `json:"file_path" jsonschema:"required,description=Path to the file within the repository"`
	Ref      string `json:"ref,omitempty" jsonschema:"description=Branch name\\, tag\\, or commit SHA. Defaults to the project's default branch."`
}

type listRepoTreeArgs struct {
	Project   string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	Path      string `json:"path,omitempty" jsonschema:"description=Directory path within the repo. Empty for root."`
	Ref       string `json:"ref,omitempty" jsonschema:"description=Branch name\\, tag\\, or commit SHA. Defaults to the project's default branch."`
	Recursive bool   `json:"recursive,omitempty" jsonschema:"description=List files recursively. Default: false"`
	PerPage   int64  `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 100)"`
}

func registerRepoFileTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_repo_file",
		Description: "Read file contents from a GitLab repository without cloning. Returns the decoded file content. Useful for AI agents that need to inspect code remotely. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getRepoFileArgs) (*mcp.CallToolResult, any, error) {
		opts := &gitlab.GetFileOptions{}
		if args.Ref != "" {
			opts.Ref = gitlab.Ptr(args.Ref)
		}

		file, _, err := client.RepositoryFiles.GetFile(args.Project, args.FilePath, opts)
		if err != nil {
			return textResult(fmt.Sprintf("Error reading file: %v", err)), nil, nil
		}

		var content string
		if file.Encoding == "base64" {
			decoded, err := base64.StdEncoding.DecodeString(file.Content)
			if err != nil {
				return textResult(fmt.Sprintf("Error decoding file: %v", err)), nil, nil
			}
			content = string(decoded)
		} else {
			content = file.Content
		}

		return textResult(content), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_repo_tree",
		Description: "List files and directories in a GitLab repository. Returns file names, types (blob/tree), and paths. Useful for understanding project structure without cloning. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listRepoTreeArgs) (*mcp.CallToolResult, any, error) {
		perPage := args.PerPage
		if perPage == 0 {
			perPage = 100
		}
		opts := &gitlab.ListTreeOptions{
			ListOptions: gitlab.ListOptions{PerPage: perPage},
		}
		if args.Path != "" {
			opts.Path = gitlab.Ptr(args.Path)
		}
		if args.Ref != "" {
			opts.Ref = gitlab.Ptr(args.Ref)
		}
		if args.Recursive {
			opts.Recursive = gitlab.Ptr(true)
		}

		tree, _, err := client.Repositories.ListTree(args.Project, opts)
		if err != nil {
			return textResult(fmt.Sprintf("Error listing tree: %v", err)), nil, nil
		}

		data, _ := json.Marshal(tree)
		return textResult(string(data)), nil, nil
	})
}
