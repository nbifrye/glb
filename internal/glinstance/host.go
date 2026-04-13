package glinstance

import (
	"fmt"
	"strings"
)

const defaultHost = "gitlab.com"

func Default() string {
	return defaultHost
}

func NormalizeHostname(hostname string) string {
	hostname = strings.TrimPrefix(hostname, "https://")
	hostname = strings.TrimPrefix(hostname, "http://")
	hostname = strings.TrimSuffix(hostname, "/")
	return strings.ToLower(hostname)
}

func APIEndpoint(hostname, protocol string) string {
	if protocol == "" {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/api/v4", protocol, NormalizeHostname(hostname))
}
