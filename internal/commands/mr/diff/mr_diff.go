package diff

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

			version, err := gitlabop.GetMergeRequestDiff(client, project, int64(mrID))
			if err != nil {
				return err
			}
			if version == nil {
				fmt.Fprintln(f.IO.Out, "No diffs found.")
				return nil
			}

			if outputJSON {
				data, err := json.MarshalIndent(version, "", "  ")
				if err != nil {
					return fmt.Errorf("marshaling response: %w", err)
				}
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
