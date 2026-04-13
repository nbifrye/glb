# glb

A GitLab CLI tool with MCP (Model Context Protocol) server support for AI agents.

## Features

### CLI Commands

`glb` provides glab-compatible commands for common GitLab operations:

- `glb auth login/status` - Authentication management
- `glb project view` - View project details
- `glb issue list/view/create/close/note` - Issue management
- `glb mr list/view/create/diff/merge/note` - Merge request management
- `glb ci list/view` - Pipeline management
- `glb api` - Raw GitLab REST API access

### MCP Server

`glb mcp serve` starts a Model Context Protocol server over stdio, exposing GitLab operations as tools for AI agents.

#### Basic Tools (glab parity)

`list_issues`, `get_issue`, `create_issue`, `close_issue`, `add_issue_note`, `list_merge_requests`, `get_merge_request`, `create_merge_request`, `get_merge_request_diff`, `merge_merge_request`, `add_mr_note`, `get_project`, `list_pipelines`, `get_pipeline`

#### Differentiated Tools (not available in glab)

| Tool | Description |
|---|---|
| `get_repo_file` | Read file contents from a repository without cloning |
| `list_repo_tree` | List files and directories in a repository |
| `search_code` | Search code across projects or within a specific project/group |
| `list_discussions` | List discussion threads on MRs/issues |
| `reply_to_discussion` | Reply to a discussion thread (not just top-level notes) |
| `get_mr_conflicts` | Get merge request conflict information |
| `list_pipeline_artifacts` | List pipeline job artifacts |
| `compare_refs` | Compare two branches, tags, or commits |
| `add_time_spent` | Add time spent on an issue/MR |
| `set_time_estimate` | Set time estimate on an issue/MR |
| `batch_update` | Bulk update multiple issues/MRs at once |

## Installation

```bash
go install github.com/nbifrye/glb/cmd/glb@latest
```

Or build from source:

```bash
make build
```

## Authentication

Set your GitLab personal access token via environment variable:

```bash
export GITLAB_TOKEN=glpat-xxxxxxxxxxxx
```

Or use the login command:

```bash
glb auth login --hostname gitlab.com --token glpat-xxxxxxxxxxxx
```

## MCP Integration

### Claude Code

Add to your MCP server configuration:

```json
{
  "mcpServers": {
    "glb": {
      "type": "stdio",
      "command": "glb",
      "args": ["mcp", "serve"]
    }
  }
}
```
