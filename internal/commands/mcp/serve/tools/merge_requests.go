package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type listMRsArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	State   string `json:"state,omitempty" jsonschema:"description=Filter by state: opened\\, closed\\, merged\\, all. Default: opened"`
	Labels  string `json:"labels,omitempty" jsonschema:"description=Comma-separated label names"`
	PerPage int64  `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 20)"`
}

type getMRArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path"`
	MRID    int64  `json:"mr_id" jsonschema:"required,description=Merge request IID"`
}

type createMRArgs struct {
	Project      string `json:"project" jsonschema:"required,description=Project path"`
	Title        string `json:"title" jsonschema:"required,description=MR title"`
	SourceBranch string `json:"source_branch" jsonschema:"required,description=Source branch"`
	TargetBranch string `json:"target_branch" jsonschema:"required,description=Target branch"`
	Description  string `json:"description,omitempty" jsonschema:"description=MR description (Markdown)"`
	Labels       string `json:"labels,omitempty" jsonschema:"description=Comma-separated label names"`
	Draft        bool   `json:"draft,omitempty" jsonschema:"description=Create as draft MR"`
}

type getMRDiffArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path"`
	MRID    int64  `json:"mr_id" jsonschema:"required,description=Merge request IID"`
}

type mergeMRArgs struct {
	Project      string `json:"project" jsonschema:"required,description=Project path"`
	MRID         int64  `json:"mr_id" jsonschema:"required,description=Merge request IID"`
	Squash       bool   `json:"squash,omitempty" jsonschema:"description=Squash commits on merge"`
	RemoveBranch bool   `json:"remove_branch,omitempty" jsonschema:"description=Remove source branch after merge"`
}

type addMRNoteArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path"`
	MRID    int64  `json:"mr_id" jsonschema:"required,description=Merge request IID"`
	Body    string `json:"body" jsonschema:"required,description=Note body (Markdown)"`
}

func registerMergeRequestTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_merge_requests",
		Description: "List merge requests in a GitLab project. Returns MR IID, title, state, source/target branches.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listMRsArgs) (*mcp.CallToolResult, any, error) {
		perPage := args.PerPage
		if perPage == 0 {
			perPage = 20
		}
		opts := &gitlab.ListProjectMergeRequestsOptions{
			ListOptions: gitlab.ListOptions{PerPage: perPage},
		}
		if args.State != "" {
			opts.State = gitlab.Ptr(args.State)
		}
		if args.Labels != "" {
			labels := splitLabels(args.Labels)
			opts.Labels = (*gitlab.LabelOptions)(&labels)
		}

		mrs, _, err := client.MergeRequests.ListProjectMergeRequests(args.Project, opts)
		if err != nil {
			return textResult(fmt.Sprintf("Error listing merge requests: %v", err)), nil, nil
		}

		data, _ := json.Marshal(mrs)
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_merge_request",
		Description: "Get detailed information about a specific merge request including description, source/target branches, and approval status.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getMRArgs) (*mcp.CallToolResult, any, error) {
		mr, _, err := client.MergeRequests.GetMergeRequest(args.Project, args.MRID, nil)
		if err != nil {
			return textResult(fmt.Sprintf("Error getting merge request: %v", err)), nil, nil
		}

		data, _ := json.Marshal(mr)
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_merge_request",
		Description: "Create a new merge request.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createMRArgs) (*mcp.CallToolResult, any, error) {
		title := args.Title
		if args.Draft {
			title = "Draft: " + title
		}

		opts := &gitlab.CreateMergeRequestOptions{
			Title:        gitlab.Ptr(title),
			SourceBranch: gitlab.Ptr(args.SourceBranch),
			TargetBranch: gitlab.Ptr(args.TargetBranch),
		}
		if args.Description != "" {
			opts.Description = gitlab.Ptr(args.Description)
		}
		if args.Labels != "" {
			labels := splitLabels(args.Labels)
			opts.Labels = (*gitlab.LabelOptions)(&labels)
		}

		mr, _, err := client.MergeRequests.CreateMergeRequest(args.Project, opts)
		if err != nil {
			return textResult(fmt.Sprintf("Error creating MR: %v", err)), nil, nil
		}

		data, _ := json.Marshal(mr)
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_merge_request_diff",
		Description: "Get the diff/changes of a merge request. Shows file-level diffs with old and new paths.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getMRDiffArgs) (*mcp.CallToolResult, any, error) {
		versions, _, err := client.MergeRequests.GetMergeRequestDiffVersions(args.Project, args.MRID, &gitlab.GetMergeRequestDiffVersionsOptions{})
		if err != nil {
			return textResult(fmt.Sprintf("Error getting diff versions: %v", err)), nil, nil
		}
		if len(versions) == 0 {
			return textResult("No diffs found."), nil, nil
		}

		version, _, err := client.MergeRequests.GetSingleMergeRequestDiffVersion(args.Project, args.MRID, versions[0].ID, &gitlab.GetSingleMergeRequestDiffVersionOptions{})
		if err != nil {
			return textResult(fmt.Sprintf("Error getting diff: %v", err)), nil, nil
		}

		data, _ := json.Marshal(version)
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "merge_merge_request",
		Description: "Merge an open merge request.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(true)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args mergeMRArgs) (*mcp.CallToolResult, any, error) {
		opts := &gitlab.AcceptMergeRequestOptions{}
		if args.Squash {
			opts.Squash = gitlab.Ptr(true)
		}
		if args.RemoveBranch {
			opts.ShouldRemoveSourceBranch = gitlab.Ptr(true)
		}

		mr, _, err := client.MergeRequests.AcceptMergeRequest(args.Project, args.MRID, opts)
		if err != nil {
			return textResult(fmt.Sprintf("Error merging MR: %v", err)), nil, nil
		}

		return textResult(fmt.Sprintf("Merged !%d: %s", mr.IID, mr.Title)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_mr_note",
		Description: "Add a comment/note to a merge request.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addMRNoteArgs) (*mcp.CallToolResult, any, error) {
		note, _, err := client.Notes.CreateMergeRequestNote(args.Project, args.MRID, &gitlab.CreateMergeRequestNoteOptions{
			Body: gitlab.Ptr(args.Body),
		})
		if err != nil {
			return textResult(fmt.Sprintf("Error adding note: %v", err)), nil, nil
		}

		return textResult(fmt.Sprintf("Added note #%d to MR !%d", note.ID, args.MRID)), nil, nil
	})
}
