package view

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "view <project>",
		Short: "View project details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			project, _, err := client.Projects.GetProject(args[0], nil)
			if err != nil {
				return fmt.Errorf("getting project: %w", err)
			}

			if outputJSON {
				data, _ := json.MarshalIndent(project, "", "  ")
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Name:        %s\n", project.NameWithNamespace)
			fmt.Fprintf(f.IO.Out, "Path:        %s\n", project.PathWithNamespace)
			fmt.Fprintf(f.IO.Out, "Description: %s\n", project.Description)
			fmt.Fprintf(f.IO.Out, "URL:         %s\n", project.WebURL)
			fmt.Fprintf(f.IO.Out, "Default branch: %s\n", project.DefaultBranch)
			fmt.Fprintf(f.IO.Out, "Visibility:  %s\n", project.Visibility)
			return nil
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")

	return cmd
}
