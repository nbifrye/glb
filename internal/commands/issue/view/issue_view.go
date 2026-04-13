package view

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/gitlabop"
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

			issue, err := gitlabop.GetIssue(client, project, int64(issueID))
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

			fmt.Fprintf(f.IO.Out, "Title:    #%d %s\n", issue.IID, issue.Title)
			fmt.Fprintf(f.IO.Out, "State:    %s\n", issue.State)
			if len(issue.Assignees) > 0 {
				names := make([]string, 0, len(issue.Assignees))
				for _, a := range issue.Assignees {
					names = append(names, a.Username)
				}
				fmt.Fprintf(f.IO.Out, "Assignee: %s\n", strings.Join(names, ", "))
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
