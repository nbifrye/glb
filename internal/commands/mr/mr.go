package mr

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	approveCmd "github.com/nbifrye/glb/internal/commands/mr/approve"
	createCmd "github.com/nbifrye/glb/internal/commands/mr/create"
	diffCmd "github.com/nbifrye/glb/internal/commands/mr/diff"
	listCmd "github.com/nbifrye/glb/internal/commands/mr/list"
	mergeCmd "github.com/nbifrye/glb/internal/commands/mr/merge"
	noteCmd "github.com/nbifrye/glb/internal/commands/mr/note"
	unapproveCmd "github.com/nbifrye/glb/internal/commands/mr/unapprove"
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
	cmd.AddCommand(approveCmd.NewCmd(f))
	cmd.AddCommand(unapproveCmd.NewCmd(f))

	return cmd
}
