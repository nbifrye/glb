package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type listIssuesArgs struct {
	Project  string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	State    string `json:"state,omitempty" jsonschema:"description=Filter by state: opened\\, closed\\, all. Default: opened"`
	Labels   string `json:"labels,omitempty" jsonschema:"description=Comma-separated label names"`
	Assignee string `json:"assignee,omitempty" jsonschema:"description=Filter by assignee username"`
	Search   string `json:"search,omitempty" jsonschema:"description=Search in title and description"`
	PerPage  int64  `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 20)"`
}

type getIssueArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path"`
	IssueID int64  `json:"issue_id" jsonschema:"required,description=Issue IID (project-level ID)"`
}

type createIssueArgs struct {
	Project     string `json:"project" jsonschema:"required,description=Project path"`
	Title       string `json:"title" jsonschema:"required,description=Issue title"`
	Description string `json:"description,omitempty" jsonschema:"description=Issue description (Markdown)"`
	Labels      string `json:"labels,omitempty" jsonschema:"description=Comma-separated label names"`
}

type closeIssueArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path"`
	IssueID int64  `json:"issue_id" jsonschema:"required,description=Issue IID"`
}

type addIssueNoteArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path"`
	IssueID int64  `json:"issue_id" jsonschema:"required,description=Issue IID"`
	Body    string `json:"body" jsonschema:"required,description=Note body (Markdown)"`
}

func registerIssueTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_issues",
		Description: "List issues in a GitLab project. Returns issue IID, title, state, labels, and assignees.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listIssuesArgs) (*mcp.CallToolResult, any, error) {
		perPage := args.PerPage
		if perPage == 0 {
			perPage = 20
		}
		opts := &gitlab.ListProjectIssuesOptions{
			ListOptions: gitlab.ListOptions{PerPage: perPage},
		}
		if args.State != "" {
			opts.State = gitlab.Ptr(args.State)
		}
		if args.Labels != "" {
			labels := splitLabels(args.Labels)
			opts.Labels = (*gitlab.LabelOptions)(&labels)
		}
		if args.Search != "" {
			opts.Search = gitlab.Ptr(args.Search)
		}
		if args.Assignee != "" {
			users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{Username: gitlab.Ptr(args.Assignee)})
			if err == nil && len(users) > 0 {
				opts.AssigneeID = gitlab.AssigneeID(users[0].ID)
			}
		}

		issues, _, err := client.Issues.ListProjectIssues(args.Project, opts)
		if err != nil {
			return textResult(fmt.Sprintf("Error listing issues: %v", err)), nil, nil
		}

		data, _ := json.Marshal(issues)
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_issue",
		Description: "Get detailed information about a specific GitLab issue including description, labels, assignees, and timestamps.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getIssueArgs) (*mcp.CallToolResult, any, error) {
		issue, _, err := client.Issues.GetIssue(args.Project, args.IssueID)
		if err != nil {
			return textResult(fmt.Sprintf("Error getting issue: %v", err)), nil, nil
		}

		data, _ := json.Marshal(issue)
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_issue",
		Description: "Create a new issue in a GitLab project.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createIssueArgs) (*mcp.CallToolResult, any, error) {
		opts := &gitlab.CreateIssueOptions{
			Title: gitlab.Ptr(args.Title),
		}
		if args.Description != "" {
			opts.Description = gitlab.Ptr(args.Description)
		}
		if args.Labels != "" {
			labels := splitLabels(args.Labels)
			opts.Labels = (*gitlab.LabelOptions)(&labels)
		}

		issue, _, err := client.Issues.CreateIssue(args.Project, opts)
		if err != nil {
			return textResult(fmt.Sprintf("Error creating issue: %v", err)), nil, nil
		}

		data, _ := json.Marshal(issue)
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "close_issue",
		Description: "Close a GitLab issue.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(true)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args closeIssueArgs) (*mcp.CallToolResult, any, error) {
		issue, _, err := client.Issues.UpdateIssue(args.Project, args.IssueID, &gitlab.UpdateIssueOptions{
			StateEvent: gitlab.Ptr("close"),
		})
		if err != nil {
			return textResult(fmt.Sprintf("Error closing issue: %v", err)), nil, nil
		}

		return textResult(fmt.Sprintf("Closed issue #%d: %s", issue.IID, issue.Title)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_issue_note",
		Description: "Add a comment/note to a GitLab issue.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addIssueNoteArgs) (*mcp.CallToolResult, any, error) {
		note, _, err := client.Notes.CreateIssueNote(args.Project, args.IssueID, &gitlab.CreateIssueNoteOptions{
			Body: gitlab.Ptr(args.Body),
		})
		if err != nil {
			return textResult(fmt.Sprintf("Error adding note: %v", err)), nil, nil
		}

		return textResult(fmt.Sprintf("Added note #%d to issue #%d", note.ID, args.IssueID)), nil, nil
	})
}
