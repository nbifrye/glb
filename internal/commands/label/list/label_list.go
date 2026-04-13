package list

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/gitlabop"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project    string
		search     string
		perPage    int
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List project labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			labels, err := gitlabop.ListLabels(client, gitlabop.ListLabelsOptions{
				Project: project,
				Search:  search,
				PerPage: int64(perPage),
			})
			if err != nil {
				return err
			}

			if outputJSON {
				data, err := json.MarshalIndent(labels, "", "  ")
				if err != nil {
					return fmt.Errorf("marshaling response: %w", err)
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(labels) == 0 {
				fmt.Fprintln(f.IO.Out, "No labels found.")
				return nil
			}

			for _, l := range labels {
				fmt.Fprintf(f.IO.Out, "%s\t%s\t(issues: %d open, %d closed)\n",
					l.Name, l.Color, l.OpenIssuesCount, l.ClosedIssuesCount)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().StringVar(&search, "search", "", "Search labels by keyword")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Number of items per page")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
