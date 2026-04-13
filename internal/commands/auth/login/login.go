package login

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/glinstance"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var hostname string
	var token string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a GitLab instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			if hostname == "" {
				hostname = glinstance.Default()
			}
			hostname = glinstance.NormalizeHostname(hostname)

			if token == "" {
				fmt.Fprintf(f.IO.Out, "Enter your GitLab personal access token for %s: ", hostname)
				reader := bufio.NewReader(f.IO.In)
				line, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("reading token: %w", err)
				}
				token = strings.TrimSpace(line)
			}

			if token == "" {
				return fmt.Errorf("token cannot be empty")
			}

			cfg.SetHost(hostname, token, "https")
			if err := cfg.Write(); err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Authenticated to %s\n", hostname)
			return nil
		},
	}

	cmd.Flags().StringVarP(&hostname, "hostname", "h", "", "GitLab hostname (default: gitlab.com)")
	cmd.Flags().StringVarP(&token, "token", "t", "", "Personal access token")

	return cmd
}
