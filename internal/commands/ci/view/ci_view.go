package view

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project    string
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "view <pipeline-id>",
		Short: "View pipeline details",
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

			pipeline, _, err := client.Pipelines.GetPipeline(project, int64(pipelineID))
			if err != nil {
				return fmt.Errorf("getting pipeline: %w", err)
			}

			if outputJSON {
				data, _ := json.MarshalIndent(pipeline, "", "  ")
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Pipeline #%d\n", pipeline.ID)
			fmt.Fprintf(f.IO.Out, "Status:  %s\n", pipeline.Status)
			fmt.Fprintf(f.IO.Out, "Ref:     %s\n", pipeline.Ref)
			fmt.Fprintf(f.IO.Out, "SHA:     %s\n", pipeline.SHA)
			fmt.Fprintf(f.IO.Out, "Created: %s\n", pipeline.CreatedAt)
			fmt.Fprintf(f.IO.Out, "URL:     %s\n", pipeline.WebURL)
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
