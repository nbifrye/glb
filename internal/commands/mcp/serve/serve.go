package serve

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory, version string) *cobra.Command {
	var hostname string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start MCP server over stdio",
		Long: `Start a Model Context Protocol (MCP) server that exposes GitLab operations
as tools for AI agents. The server communicates via JSON-RPC over stdio.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(cmd.Context(), f, version, hostname)
		},
	}

	cmd.Flags().StringVar(&hostname, "hostname", "", "GitLab hostname to connect to (overrides config)")

	return cmd
}
