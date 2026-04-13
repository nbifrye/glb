package serve

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/commands/mcp/serve/tools"
)

func runServer(ctx context.Context, f *cmdutils.Factory, version, hostname string) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "glb",
		Version: version,
	}, &mcp.ServerOptions{
		Instructions: "glb is a GitLab MCP server. Use project paths in the format 'group/project' (e.g., 'gitlab-org/gitlab'). " +
			"Resource IDs (issue_id, mr_id) are project-scoped IIDs, not global IDs. " +
			"State filters accept: 'opened', 'closed', 'all' for issues; 'opened', 'closed', 'merged', 'all' for merge requests. " +
			"Tools marked as read-only are safe to call without side effects.",
	})

	var (
		gitlabClient *gitlab.Client
		err          error
	)
	if hostname != "" {
		gitlabClient, err = f.GitLabClientForHost(hostname)
	} else {
		gitlabClient, err = f.GitLabClient()
	}
	if err != nil {
		return fmt.Errorf("initializing GitLab client: %w", err)
	}

	tools.RegisterAll(server, gitlabClient)

	return server.Run(ctx, &mcp.StdioTransport{})
}
