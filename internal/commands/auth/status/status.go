package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/auth"
	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			hostname, token := auth.DefaultHostWithToken(cfg)
			if token == "" {
				fmt.Fprintln(f.IO.Out, "Not authenticated. Run 'glb auth login' or set GITLAB_TOKEN.")
				return nil
			}

			var masked string
			switch {
			case len(token) <= 4:
				masked = "****"
			case len(token) < 8:
				masked = token[:2] + "****"
			default:
				masked = token[:4] + "****" + token[len(token)-4:]
			}
			fmt.Fprintf(f.IO.Out, "Authenticated to %s (token: %s)\n", hostname, masked)
			return nil
		},
	}
}
