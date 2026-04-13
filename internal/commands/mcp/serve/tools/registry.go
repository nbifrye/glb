package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func RegisterAll(s *mcp.Server, client *gitlab.Client) {
	registerIssueTools(s, client)
	registerMergeRequestTools(s, client)
	registerProjectTools(s, client)
	registerPipelineTools(s, client)
	registerRepoFileTools(s, client)
	registerSearchTools(s, client)
	registerDiscussionTools(s, client)
	registerMRConflictTools(s, client)
	registerArtifactTools(s, client)
	registerCompareTools(s, client)
	registerTimetrackingTools(s, client)
	registerBatchTools(s, client)
	registerLabelTools(s, client)
}
