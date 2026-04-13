package project

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	viewCmd "github.com/nbifrye/glb/internal/commands/project/view"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project <command>",
		Short:   "Manage GitLab projects",
		Aliases: []string{"repo"},
	}

	cmd.AddCommand(viewCmd.NewCmd(f))

	return cmd
}
