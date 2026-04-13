package jobs

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/gitlabop"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project    string
		scope      string
		perPage    int
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "jobs <pipeline-id>",
		Short: "List jobs in a pipeline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pipelineID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid pipeline ID: %s", args[0])
			}

			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			jobs, err := gitlabop.ListPipelineJobs(client, gitlabop.ListPipelineJobsOptions{
				Project:    project,
				PipelineID: int64(pipelineID),
				Scope:      scope,
				PerPage:    int64(perPage),
			})
			if err != nil {
				return err
			}

			if outputJSON {
				data, err := json.MarshalIndent(jobs, "", "  ")
				if err != nil {
					return fmt.Errorf("marshaling response: %w", err)
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(jobs) == 0 {
				fmt.Fprintln(f.IO.Out, "No jobs found.")
				return nil
			}

			for _, j := range jobs {
				fmt.Fprintf(f.IO.Out, "#%d\t%s\t%s\t%s\n", j.ID, j.Name, j.Stage, j.Status)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().StringVar(&scope, "scope", "", "Filter by scope: created, pending, running, failed, success, canceled, skipped")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Number of items per page")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
