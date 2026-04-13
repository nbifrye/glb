package api

import (
	"fmt"

	"github.com/nbifrye/glb/internal/glinstance"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func NewGitLabClient(token, hostname, apiProtocol string) (*gitlab.Client, error) {
	baseURL := glinstance.APIEndpoint(hostname, apiProtocol)
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, fmt.Errorf("creating GitLab client: %w", err)
	}
	return client, nil
}
