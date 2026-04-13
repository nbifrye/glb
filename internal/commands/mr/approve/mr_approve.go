package approve

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/gitlabop"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var project string

	cmd := &cobra.Command{
		Use:   "approve <mr-iid>",
		Short: "Approve a merge request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mrIID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid merge request IID: %s", args[0])
			}

			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			_, err = gitlabop.ApproveMergeRequest(client, project, int64(mrIID))
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Approved merge request !%d\n", mrIID)
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
