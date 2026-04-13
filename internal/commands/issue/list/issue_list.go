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
		assignee   string
		perPage    int
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			issues, err := gitlabop.ListIssues(client, gitlabop.ListIssuesOptions{
				Project:  project,
				State:    state,
				Labels:   labels,
				Assignee: assignee,
				PerPage:  int64(perPage),
			})
			if err != nil {
				return err
			}

			if outputJSON {
				data, err := json.MarshalIndent(issues, "", "  ")
				if err != nil {
					return fmt.Errorf("marshaling response: %w", err)
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(issues) == 0 {
				fmt.Fprintln(f.IO.Out, "No issues found.")
				return nil
			}

			for _, issue := range issues {
				fmt.Fprintf(f.IO.Out, "#%d\t%s\t(%s)\n", issue.IID, issue.Title, issue.State)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().StringVarP(&state, "state", "s", "opened", "Filter by state: opened, closed, all")
	cmd.Flags().StringSliceVarP(&labels, "labels", "l", nil, "Filter by labels")
	cmd.Flags().StringVarP(&assignee, "assignee", "a", "", "Filter by assignee username")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Number of items per page")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
