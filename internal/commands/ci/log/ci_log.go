package log

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
		Use:   "log <job-id>",
		Short: "View job log output",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jobID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid job ID: %s", args[0])
			}

			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			logOutput, err := gitlabop.GetJobLog(client, project, int64(jobID))
			if err != nil {
				return err
			}

			fmt.Fprint(f.IO.Out, logOutput)
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
