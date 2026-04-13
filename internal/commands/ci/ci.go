package ci

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	jobsCmd "github.com/nbifrye/glb/internal/commands/ci/jobs"
	listCmd "github.com/nbifrye/glb/internal/commands/ci/list"
	logCmd "github.com/nbifrye/glb/internal/commands/ci/log"
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
	cmd.AddCommand(jobsCmd.NewCmd(f))
	cmd.AddCommand(logCmd.NewCmd(f))

	return cmd
}
