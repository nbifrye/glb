package tools

import (
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

func errorResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
		IsError: true,
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func clampPerPage(perPage, defaultVal int64) int64 {
	if perPage <= 0 {
		return defaultVal
	}
	if perPage > 100 {
		return 100
	}
	return perPage
}

func splitLabels(s string) []string {
	parts := strings.Split(s, ",")
	labels := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			labels = append(labels, p)
		}
	}
	return labels
}
