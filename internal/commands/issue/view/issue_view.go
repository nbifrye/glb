package view

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project    string
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "view <issue-id>",
		Short: "View an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[0])
			}

			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			issue, _, err := client.Issues.GetIssue(project, int64(issueID))
			if err != nil {
				return fmt.Errorf("getting issue: %w", err)
			}

			if outputJSON {
				data, _ := json.MarshalIndent(issue, "", "  ")
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Title:    #%d %s\n", issue.IID, issue.Title)
			fmt.Fprintf(f.IO.Out, "State:    %s\n", issue.State)
			if issue.Assignee != nil {
				fmt.Fprintf(f.IO.Out, "Assignee: %s\n", issue.Assignee.Username)
			}
			if len(issue.Labels) > 0 {
				fmt.Fprintf(f.IO.Out, "Labels:   %s\n", strings.Join(issue.Labels, ", "))
			}
			fmt.Fprintf(f.IO.Out, "URL:      %s\n", issue.WebURL)
			if issue.Description != "" {
				fmt.Fprintf(f.IO.Out, "\n%s\n", issue.Description)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
