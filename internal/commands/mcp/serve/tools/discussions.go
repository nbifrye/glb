package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type listDiscussionsArgs struct {
	Project      string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	ResourceType string `json:"resource_type" jsonschema:"required,description=Resource type: 'mr' or 'issue'"`
	ResourceID   int64  `json:"resource_id" jsonschema:"required,description=Resource IID (merge request or issue IID)"`
	PerPage      int64  `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 20)"`
}

type replyToDiscussionArgs struct {
	Project      string `json:"project" jsonschema:"required,description=Project path"`
	ResourceType string `json:"resource_type" jsonschema:"required,description=Resource type: 'mr' or 'issue'"`
	ResourceID   int64  `json:"resource_id" jsonschema:"required,description=Resource IID"`
	DiscussionID string `json:"discussion_id" jsonschema:"required,description=Discussion ID"`
	Body         string `json:"body" jsonschema:"required,description=Reply text (Markdown)"`
}

func registerDiscussionTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_discussions",
		Description: "List discussion threads on a merge request or issue. Returns full threaded conversations including all notes/replies. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listDiscussionsArgs) (*mcp.CallToolResult, any, error) {
		perPage := args.PerPage
		if perPage == 0 {
			perPage = 20
		}

		switch args.ResourceType {
		case "mr":
			discussions, _, err := client.Discussions.ListMergeRequestDiscussions(args.Project, args.ResourceID, &gitlab.ListMergeRequestDiscussionsOptions{
				ListOptions: gitlab.ListOptions{PerPage: perPage},
			})
			if err != nil {
				return errorResult(fmt.Sprintf("Error listing discussions: %v", err)), nil, nil
			}
			data, err := json.Marshal(discussions)
			if err != nil {
				return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
			}
			return textResult(string(data)), nil, nil

		case "issue":
			discussions, _, err := client.Discussions.ListIssueDiscussions(args.Project, args.ResourceID, &gitlab.ListIssueDiscussionsOptions{
				ListOptions: gitlab.ListOptions{PerPage: perPage},
			})
			if err != nil {
				return errorResult(fmt.Sprintf("Error listing discussions: %v", err)), nil, nil
			}
			data, err := json.Marshal(discussions)
			if err != nil {
				return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
			}
			return textResult(string(data)), nil, nil

		default:
			return errorResult("Error: resource_type must be 'mr' or 'issue'"), nil, nil
		}
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "reply_to_discussion",
		Description: "Reply to a discussion thread on a merge request or issue. Unlike glab's note command which only creates top-level comments, this replies within an existing thread. [DIFFERENTIATED: not available in glab]",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(false)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args replyToDiscussionArgs) (*mcp.CallToolResult, any, error) {
		switch args.ResourceType {
		case "mr":
			note, _, err := client.Discussions.AddMergeRequestDiscussionNote(args.Project, args.ResourceID, args.DiscussionID, &gitlab.AddMergeRequestDiscussionNoteOptions{
				Body: gitlab.Ptr(args.Body),
			})
			if err != nil {
				return errorResult(fmt.Sprintf("Error replying: %v", err)), nil, nil
			}
			return textResult(fmt.Sprintf("Added reply (note #%d) to discussion %s", note.ID, args.DiscussionID)), nil, nil

		case "issue":
			note, _, err := client.Discussions.AddIssueDiscussionNote(args.Project, args.ResourceID, args.DiscussionID, &gitlab.AddIssueDiscussionNoteOptions{
				Body: gitlab.Ptr(args.Body),
			})
			if err != nil {
				return errorResult(fmt.Sprintf("Error replying: %v", err)), nil, nil
			}
			return textResult(fmt.Sprintf("Added reply (note #%d) to discussion %s", note.ID, args.DiscussionID)), nil, nil

		default:
			return errorResult("Error: resource_type must be 'mr' or 'issue'"), nil, nil
		}
	})
}
