package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/gitlabop"
)

type listMRsArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path (e.g. 'group/project')"`
	State   string `json:"state,omitempty" jsonschema:"description=Filter by state: opened\\, closed\\, merged\\, all. Default: opened"`
	Labels  string `json:"labels,omitempty" jsonschema:"description=Comma-separated label names"`
	PerPage int64  `json:"per_page,omitempty" jsonschema:"description=Results per page (default: 20\\, max: 100)"`
	Page    int64  `json:"page,omitempty" jsonschema:"description=Page number for pagination (default: 1)"`
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

type approveMRArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path"`
	MRID    int64  `json:"mr_id" jsonschema:"required,description=Merge request IID"`
}

type unapproveMRArgs struct {
	Project string `json:"project" jsonschema:"required,description=Project path"`
	MRID    int64  `json:"mr_id" jsonschema:"required,description=Merge request IID"`
}

func registerMergeRequestTools(s *mcp.Server, client *gitlab.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_merge_requests",
		Description: "List merge requests in a GitLab project. Returns MR IID, title, state, source/target branches.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listMRsArgs) (*mcp.CallToolResult, any, error) {
		var labels []string
		if args.Labels != "" {
			labels = splitLabels(args.Labels)
		}

		mrs, err := gitlabop.ListMergeRequests(client, gitlabop.ListMergeRequestsOptions{
			Project: args.Project,
			State:   args.State,
			Labels:  labels,
			PerPage: clampPerPage(args.PerPage, 20),
			Page:    args.Page,
		})
		if err != nil {
			return errorResult(fmt.Sprintf("Error listing merge requests: %v", err)), nil, nil
		}

		data, err := json.Marshal(mrs)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_merge_request",
		Description: "Get detailed information about a specific merge request including description, source/target branches, and approval status.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getMRArgs) (*mcp.CallToolResult, any, error) {
		mr, err := gitlabop.GetMergeRequest(client, args.Project, args.MRID)
		if err != nil {
			return errorResult(fmt.Sprintf("Error getting merge request: %v", err)), nil, nil
		}

		data, err := json.Marshal(mr)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_merge_request",
		Description: "Create a new merge request.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(false)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createMRArgs) (*mcp.CallToolResult, any, error) {
		var labels []string
		if args.Labels != "" {
			labels = splitLabels(args.Labels)
		}

		mr, err := gitlabop.CreateMergeRequest(client, gitlabop.CreateMergeRequestOptions{
			Project:      args.Project,
			Title:        args.Title,
			SourceBranch: args.SourceBranch,
			TargetBranch: args.TargetBranch,
			Description:  args.Description,
			Labels:       labels,
			Draft:        args.Draft,
		})
		if err != nil {
			return errorResult(fmt.Sprintf("Error creating MR: %v", err)), nil, nil
		}

		data, err := json.Marshal(mr)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_merge_request_diff",
		Description: "Get the diff/changes of a merge request. Shows file-level diffs with old and new paths.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getMRDiffArgs) (*mcp.CallToolResult, any, error) {
		version, err := gitlabop.GetMergeRequestDiff(client, args.Project, args.MRID)
		if err != nil {
			return errorResult(fmt.Sprintf("Error getting diff: %v", err)), nil, nil
		}
		if version == nil {
			return textResult("No diffs found."), nil, nil
		}

		data, err := json.Marshal(version)
		if err != nil {
			return errorResult(fmt.Sprintf("Error marshaling response: %v", err)), nil, nil
		}
		return textResult(string(data)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "merge_merge_request",
		Description: "Merge an open merge request.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(true)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args mergeMRArgs) (*mcp.CallToolResult, any, error) {
		mr, err := gitlabop.MergeMergeRequest(client, gitlabop.MergeMergeRequestOptions{
			Project:      args.Project,
			IID:          args.MRID,
			Squash:       args.Squash,
			RemoveBranch: args.RemoveBranch,
		})
		if err != nil {
			return errorResult(fmt.Sprintf("Error merging MR: %v", err)), nil, nil
		}

		return textResult(fmt.Sprintf("Merged !%d: %s", mr.IID, mr.Title)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_mr_note",
		Description: "Add a comment/note to a merge request.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(false)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addMRNoteArgs) (*mcp.CallToolResult, any, error) {
		note, err := gitlabop.AddMergeRequestNote(client, args.Project, args.MRID, args.Body)
		if err != nil {
			return errorResult(fmt.Sprintf("Error adding note: %v", err)), nil, nil
		}

		return textResult(fmt.Sprintf("Added note #%d to MR !%d", note.ID, args.MRID)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "approve_merge_request",
		Description: "Approve a merge request.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(false)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args approveMRArgs) (*mcp.CallToolResult, any, error) {
		_, err := gitlabop.ApproveMergeRequest(client, args.Project, args.MRID)
		if err != nil {
			return errorResult(fmt.Sprintf("Error approving MR: %v", err)), nil, nil
		}

		return textResult(fmt.Sprintf("Approved MR !%d", args.MRID)), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "unapprove_merge_request",
		Description: "Unapprove (revoke approval of) a merge request.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(true)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args unapproveMRArgs) (*mcp.CallToolResult, any, error) {
		err := gitlabop.UnapproveMergeRequest(client, args.Project, args.MRID)
		if err != nil {
			return errorResult(fmt.Sprintf("Error unapproving MR: %v", err)), nil, nil
		}

		return textResult(fmt.Sprintf("Unapproved MR !%d", args.MRID)), nil, nil
	})
}
