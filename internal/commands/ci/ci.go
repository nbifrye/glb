package ci

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	listCmd "github.com/nbifrye/glb/internal/commands/ci/list"
	viewCmd "github.com/nbifrye/glb/internal/commands/ci/view"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ci <command>",
		Short:   "Manage CI/CD pipelines",
		Aliases: []string{"pipeline"},
	}

	cmd.AddCommand(listCmd.NewCmd(f))
	cmd.AddCommand(viewCmd.NewCmd(f))

	return cmd
}
