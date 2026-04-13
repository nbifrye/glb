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
		Use:   "note <issue-id>",
		Short: "Add a note to an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid issue ID: %s", args[0])
			}

			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			note, _, err := client.Notes.CreateIssueNote(project, int64(issueID), &gitlab.CreateIssueNoteOptions{
				Body: gitlab.Ptr(body),
			})
			if err != nil {
				return fmt.Errorf("adding note: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "Added note #%d to issue #%d\n", note.ID, issueID)
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().StringVarP(&body, "body", "b", "", "Note body (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("body")

	return cmd
}
