package gitlabop

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func GetProject(client *gitlab.Client, project string) (*gitlab.Project, error) {
	p, _, err := client.Projects.GetProject(project, nil)
	if err != nil {
		return nil, fmt.Errorf("getting project: %w", err)
	}
	return p, nil
}
