package create

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/gitlabop"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project     string
		title       string
		description string
		labels      []string
		assignees   []int64
		outputJSON  bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			issue, err := gitlabop.CreateIssue(client, gitlabop.CreateIssueOptions{
				Project:     project,
				Title:       title,
				Description: description,
				Labels:      labels,
				AssigneeIDs: assignees,
			})
			if err != nil {
				return err
			}

			if outputJSON {
				data, err := json.MarshalIndent(issue, "", "  ")
				if err != nil {
					return fmt.Errorf("marshaling response: %w", err)
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Created issue #%d: %s\n%s\n", issue.IID, issue.Title, issue.WebURL)
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().StringVarP(&title, "title", "t", "", "Issue title (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Issue description")
	cmd.Flags().StringSliceVarP(&labels, "labels", "l", nil, "Labels")
	cmd.Flags().Int64SliceVar(&assignees, "assignees", nil, "Assignee user IDs")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("title")

	return cmd
}
