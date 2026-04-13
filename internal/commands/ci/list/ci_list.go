package list

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project    string
		status     string
		ref        string
		perPage    int
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pipelines",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			opts := &gitlab.ListProjectPipelinesOptions{
				ListOptions: gitlab.ListOptions{PerPage: int64(perPage)},
			}
			if status != "" {
				s := gitlab.BuildStateValue(status)
				opts.Status = &s
			}
			if ref != "" {
				opts.Ref = gitlab.Ptr(ref)
			}

			pipelines, _, err := client.Pipelines.ListProjectPipelines(project, opts)
			if err != nil {
				return fmt.Errorf("listing pipelines: %w", err)
			}

			if outputJSON {
				data, _ := json.MarshalIndent(pipelines, "", "  ")
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			if len(pipelines) == 0 {
				fmt.Fprintln(f.IO.Out, "No pipelines found.")
				return nil
			}

			for _, p := range pipelines {
				fmt.Fprintf(f.IO.Out, "#%d\t%s\t%s\t(%s)\n", p.ID, p.Ref, p.Status, p.CreatedAt)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status: running, pending, success, failed, etc.")
	cmd.Flags().StringVar(&ref, "ref", "", "Filter by ref (branch or tag)")
	cmd.Flags().IntVar(&perPage, "per-page", 30, "Number of items per page")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
