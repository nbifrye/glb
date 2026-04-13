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
		status     string
		ref        string
		perPage    int
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pipelines",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			pipelines, err := gitlabop.ListPipelines(client, gitlabop.ListPipelinesOptions{
				Project: project,
				Status:  status,
				Ref:     ref,
				PerPage: int64(perPage),
			})
			if err != nil {
				return err
			}

			if outputJSON {
				data, err := json.MarshalIndent(pipelines, "", "  ")
				if err != nil {
					return fmt.Errorf("marshaling response: %w", err)
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(pipelines) == 0 {
				fmt.Fprintln(f.IO.Out, "No pipelines found.")
				return nil
			}

			for _, p := range pipelines {
				fmt.Fprintf(f.IO.Out, "#%d\t%s\t%s\t(%s)\n", p.ID, p.Ref, p.Status, p.CreatedAt)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status: running, pending, success, failed, etc.")
	cmd.Flags().StringVar(&ref, "ref", "", "Filter by ref (branch or tag)")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Number of items per page")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
