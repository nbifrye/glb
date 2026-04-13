package label

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	listCmd "github.com/nbifrye/glb/internal/commands/label/list"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label <command>",
		Short: "Manage labels",
	}

	cmd.AddCommand(listCmd.NewCmd(f))

	return cmd
}
