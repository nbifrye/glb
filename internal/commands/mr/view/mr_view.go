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
		Use:   "view <mr-id>",
		Short: "View a merge request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mrID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid MR ID: %s", args[0])
			}

			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			mr, err := gitlabop.GetMergeRequest(client, project, int64(mrID))
			if err != nil {
				return err
			}

			if outputJSON {
				data, err := json.MarshalIndent(mr, "", "  ")
				if err != nil {
					return fmt.Errorf("marshaling response: %w", err)
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Title:    !%d %s\n", mr.IID, mr.Title)
			fmt.Fprintf(f.IO.Out, "State:    %s\n", mr.State)
			fmt.Fprintf(f.IO.Out, "Branch:   %s -> %s\n", mr.SourceBranch, mr.TargetBranch)
			if mr.Author != nil {
				fmt.Fprintf(f.IO.Out, "Author:   %s\n", mr.Author.Username)
			}
			if len(mr.Labels) > 0 {
				fmt.Fprintf(f.IO.Out, "Labels:   %s\n", strings.Join(mr.Labels, ", "))
			}
			fmt.Fprintf(f.IO.Out, "URL:      %s\n", mr.WebURL)
			if mr.Description != "" {
				fmt.Fprintf(f.IO.Out, "\n%s\n", mr.Description)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
