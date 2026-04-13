package serve

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/commands/mcp/serve/tools"
)

func runServer(ctx context.Context, f *cmdutils.Factory) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "glb",
		Version: "0.1.0",
	}, nil)

	gitlabClient, err := f.GitLabClient()
	if err != nil {
		return fmt.Errorf("initializing GitLab client: %w", err)
	}

	tools.RegisterAll(server, gitlabClient)

	return server.Run(ctx, &mcp.StdioTransport{})
}
