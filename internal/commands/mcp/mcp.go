package mcp

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	serveCmd "github.com/nbifrye/glb/internal/commands/mcp/serve"
)

func NewCmd(f *cmdutils.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp <command>",
		Short: "Model Context Protocol server for AI agents",
	}

	cmd.AddCommand(serveCmd.NewCmd(f, version))

	return cmd
}
