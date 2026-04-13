package list

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/cmdutils"
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

			opts := &gitlab.ListProjectIssuesOptions{
				ListOptions: gitlab.ListOptions{PerPage: int64(perPage)},
			}
			if state != "" {
				opts.State = gitlab.Ptr(state)
			}
			if len(labels) > 0 {
				opts.Labels = (*gitlab.LabelOptions)(&labels)
			}
			if assignee != "" {
				users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{Username: gitlab.Ptr(assignee)})
				if err == nil && len(users) > 0 {
					opts.AssigneeID = gitlab.AssigneeID(users[0].ID)
				}
			}

			issues, _, err := client.Issues.ListProjectIssues(project, opts)
			if err != nil {
				return fmt.Errorf("listing issues: %w", err)
			}

			if outputJSON {
				data, _ := json.MarshalIndent(issues, "", "  ")
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
	cmd.Flags().IntVar(&perPage, "per-page", 30, "Number of items per page")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
