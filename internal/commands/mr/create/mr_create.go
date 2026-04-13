package create

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/cmdutils"
	"github.com/nbifrye/glb/internal/gitlabop"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		project      string
		title        string
		description  string
		sourceBranch string
		targetBranch string
		labels       []string
		draft        bool
		outputJSON   bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a merge request",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.GitLabClient()
			if err != nil {
				return err
			}

			mr, err := gitlabop.CreateMergeRequest(client, gitlabop.CreateMergeRequestOptions{
				Project:      project,
				Title:        title,
				SourceBranch: sourceBranch,
				TargetBranch: targetBranch,
				Description:  description,
				Labels:       labels,
				Draft:        draft,
			})
			if err != nil {
				return err
			}

			if outputJSON {
				data, err := json.MarshalIndent(mr, "", "  ")
				if err != nil {
					return fmt.Errorf("marshaling response: %w", err)
				}
				fmt.Fprintln(f.IO.Out, string(data))
				return nil
			}

			fmt.Fprintf(f.IO.Out, "Created merge request !%d: %s\n%s\n", mr.IID, mr.Title, mr.WebURL)
			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project path (required)")
	cmd.Flags().StringVarP(&title, "title", "t", "", "MR title (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "MR description")
	cmd.Flags().StringVarP(&sourceBranch, "source", "s", "", "Source branch (required)")
	cmd.Flags().StringVar(&targetBranch, "target", "", "Target branch (required)")
	cmd.Flags().StringSliceVarP(&labels, "labels", "l", nil, "Labels")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create as draft MR")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("source")
	_ = cmd.MarkFlagRequired("target")

	return cmd
}
