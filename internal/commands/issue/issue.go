package issue

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	closeCmd "github.com/nbifrye/glb/internal/commands/issue/close"
	createCmd "github.com/nbifrye/glb/internal/commands/issue/create"
	listCmd "github.com/nbifrye/glb/internal/commands/issue/list"
	noteCmd "github.com/nbifrye/glb/internal/commands/issue/note"
	viewCmd "github.com/nbifrye/glb/internal/commands/issue/view"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue <command>",
		Short: "Manage GitLab issues",
	}

	cmd.AddCommand(listCmd.NewCmd(f))
	cmd.AddCommand(viewCmd.NewCmd(f))
	cmd.AddCommand(createCmd.NewCmd(f))
	cmd.AddCommand(closeCmd.NewCmd(f))
	cmd.AddCommand(noteCmd.NewCmd(f))

	return cmd
}
