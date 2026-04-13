package create

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/nbifrye/glb/internal/cmdutils"
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

			opts := &gitlab.CreateMergeRequestOptions{
				Title:        gitlab.Ptr(title),
				SourceBranch: gitlab.Ptr(sourceBranch),
				TargetBranch: gitlab.Ptr(targetBranch),
			}
			if description != "" {
				opts.Description = gitlab.Ptr(description)
			}
			if len(labels) > 0 {
				opts.Labels = (*gitlab.LabelOptions)(&labels)
			}
			if draft {
				t := "Draft: " + title
				opts.Title = gitlab.Ptr(t)
			}

			mr, _, err := client.MergeRequests.CreateMergeRequest(project, opts)
			if err != nil {
				return fmt.Errorf("creating merge request: %w", err)
			}

			if outputJSON {
				data, _ := json.MarshalIndent(mr, "", "  ")
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
