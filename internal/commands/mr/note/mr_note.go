package note

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project string
		body    string
	)

	cmd := &cobra.Command{
		Use:   "note <mr-id>",
		Short: "Add a note to a merge request",
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

			note, _, err := client.Notes.CreateMergeRequestNote(project, int64(mrID), &gitlab.CreateMergeRequestNoteOptions{
				Body: gitlab.Ptr(body),
			})
			if err != nil {
				return fmt.Errorf("adding note: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "Added note #%d to MR !%d\n", note.ID, mrID)
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().StringVarP(&body, "body", "b", "", "Note body (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("body")

	return cmd
}
