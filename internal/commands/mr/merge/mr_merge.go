package merge

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project    string
		squash     bool
		removeBranch bool
	)

	cmd := &cobra.Command{
		Use:   "merge <mr-id>",
		Short: "Merge a merge request",
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

			opts := &gitlab.AcceptMergeRequestOptions{}
			if squash {
				opts.Squash = gitlab.Ptr(true)
			}
			if removeBranch {
				opts.ShouldRemoveSourceBranch = gitlab.Ptr(true)
			}

			mr, _, err := client.MergeRequests.AcceptMergeRequest(project, int64(mrID), opts)
			if err != nil {
				return fmt.Errorf("merging MR: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "Merged !%d: %s\n", mr.IID, mr.Title)
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().BoolVar(&squash, "squash", false, "Squash commits")
	cmd.Flags().BoolVar(&removeBranch, "remove-branch", false, "Remove source branch after merge")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
