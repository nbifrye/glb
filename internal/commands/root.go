package commands

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	apiCmd "github.com/nbifrye/glb/internal/commands/api"
	authCmd "github.com/nbifrye/glb/internal/commands/auth"
	ciCmd "github.com/nbifrye/glb/internal/commands/ci"
	issueCmd "github.com/nbifrye/glb/internal/commands/issue"
	labelCmd "github.com/nbifrye/glb/internal/commands/label"
	mcpCmd "github.com/nbifrye/glb/internal/commands/mcp"
	mrCmd "github.com/nbifrye/glb/internal/commands/mr"
	projectCmd "github.com/nbifrye/glb/internal/commands/project"
	versionCmd "github.com/nbifrye/glb/internal/commands/version"
)

func NewRootCmd(f *cmdutils.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "glb <command> <subcommand> [flags]",
		Short:         "A GitLab CLI tool with MCP server support",
		Long:          "glb is a GitLab CLI that exposes GitLab operations as both CLI commands and MCP tools for AI agents.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(versionCmd.NewCmd(version))
	cmd.AddCommand(authCmd.NewCmd(f))
	cmd.AddCommand(projectCmd.NewCmd(f))
	cmd.AddCommand(issueCmd.NewCmd(f))
	cmd.AddCommand(mrCmd.NewCmd(f))
	cmd.AddCommand(ciCmd.NewCmd(f))
	cmd.AddCommand(labelCmd.NewCmd(f))
	cmd.AddCommand(apiCmd.NewCmd(f))
	cmd.AddCommand(mcpCmd.NewCmd(f, version))

	return cmd
}
