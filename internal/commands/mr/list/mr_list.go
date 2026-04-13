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
		state      string
		labels     []string
		perPage    int
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List merge requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			mrs, err := gitlabop.ListMergeRequests(client, gitlabop.ListMergeRequestsOptions{
				Project: project,
				State:   state,
				Labels:  labels,
				PerPage: int64(perPage),
			})
			if err != nil {
				return err
			}

			if outputJSON {
				data, err := json.MarshalIndent(mrs, "", "  ")
				if err != nil {
					return fmt.Errorf("marshaling response: %w", err)
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(mrs) == 0 {
				fmt.Fprintln(f.IO.Out, "No merge requests found.")
				return nil
			}

			for _, mr := range mrs {
				fmt.Fprintf(f.IO.Out, "!%d\t%s\t(%s)\n", mr.IID, mr.Title, mr.State)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().StringVarP(&state, "state", "s", "opened", "Filter by state: opened, closed, merged, all")
	cmd.Flags().StringSliceVarP(&labels, "labels", "l", nil, "Filter by labels")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Number of items per page")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
