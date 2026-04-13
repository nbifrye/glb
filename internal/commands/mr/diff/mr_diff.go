package diff

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project    string
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "diff <mr-id>",
		Short: "View merge request diff",
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

			versions, _, err := client.MergeRequests.GetMergeRequestDiffVersions(project, int64(mrID), &gitlab.GetMergeRequestDiffVersionsOptions{})
			if err != nil {
				return fmt.Errorf("getting MR diff versions: %w", err)
			}

			if len(versions) == 0 {
				fmt.Fprintln(f.IO.Out, "No diffs found.")
				return nil
			}

			latest := versions[0]
			version, _, err := client.MergeRequests.GetSingleMergeRequestDiffVersion(project, int64(mrID), latest.ID, &gitlab.GetSingleMergeRequestDiffVersionOptions{})
			if err != nil {
				return fmt.Errorf("getting diff version: %w", err)
			}

			if outputJSON {
				data, _ := json.MarshalIndent(version, "", "  ")
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			for _, d := range version.Diffs {
				fmt.Fprintf(f.IO.Out, "--- %s\n+++ %s\n", d.OldPath, d.NewPath)
				fmt.Fprintln(f.IO.Out, d.Diff)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
