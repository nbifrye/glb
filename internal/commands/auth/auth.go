package auth

import (
	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	loginCmd "github.com/nbifrye/glb/internal/commands/auth/login"
	statusCmd "github.com/nbifrye/glb/internal/commands/auth/status"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Manage glb authentication",
	}

	cmd.AddCommand(loginCmd.NewCmd(f))
	cmd.AddCommand(statusCmd.NewCmd(f))

	return cmd
}
