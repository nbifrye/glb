package merge

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/gitlabop"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project      string
		squash       bool
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

			mr, err := gitlabop.MergeMergeRequest(client, gitlabop.MergeMergeRequestOptions{
				Project:      project,
				IID:          int64(mrID),
				Squash:       squash,
				RemoveBranch: removeBranch,
			})
			if err != nil {
				return err
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
