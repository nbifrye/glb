package close

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var project string

	cmd := &cobra.Command{
		Use:   "close <issue-id>",
		Short: "Close an issue",
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

			_, _, err = client.Issues.UpdateIssue(project, int64(issueID), &gitlab.UpdateIssueOptions{
				StateEvent: gitlab.Ptr("close"),
			})
			if err != nil {
				return fmt.Errorf("closing issue: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "Closed issue #%d\n", issueID)
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
