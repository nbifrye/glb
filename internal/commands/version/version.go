package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show glb version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "glb version %s\n", version)
		},
	}
}
