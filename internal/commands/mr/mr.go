package mr

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	createCmd "github.com/nbifrye/glb/internal/commands/mr/create"
	diffCmd "github.com/nbifrye/glb/internal/commands/mr/diff"
	listCmd "github.com/nbifrye/glb/internal/commands/mr/list"
	mergeCmd "github.com/nbifrye/glb/internal/commands/mr/merge"
	noteCmd "github.com/nbifrye/glb/internal/commands/mr/note"
	viewCmd "github.com/nbifrye/glb/internal/commands/mr/view"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mr <command>",
		Short: "Manage merge requests",
	}

	cmd.AddCommand(listCmd.NewCmd(f))
	cmd.AddCommand(viewCmd.NewCmd(f))
	cmd.AddCommand(createCmd.NewCmd(f))
	cmd.AddCommand(diffCmd.NewCmd(f))
	cmd.AddCommand(mergeCmd.NewCmd(f))
	cmd.AddCommand(noteCmd.NewCmd(f))

	return cmd
}
