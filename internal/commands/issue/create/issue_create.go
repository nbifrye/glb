package create

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/cmdutils"
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

			opts := &gitlab.CreateIssueOptions{
				Title: gitlab.Ptr(title),
			}
			if description != "" {
				opts.Description = gitlab.Ptr(description)
			}
			if len(labels) > 0 {
				opts.Labels = (*gitlab.LabelOptions)(&labels)
			}
			if len(assignees) > 0 {
				opts.AssigneeIDs = &assignees
			}

			issue, _, err := client.Issues.CreateIssue(project, opts)
			if err != nil {
				return fmt.Errorf("creating issue: %w", err)
			}

			if outputJSON {
				data, _ := json.MarshalIndent(issue, "", "  ")
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
