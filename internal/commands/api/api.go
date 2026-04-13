package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/nbifrye/glb/internal/auth"
	"github.com/nbifrye/glb/internal/cmdutils"
)

func NewCmd(f *cmdutils.Factory) *cobra.Command {
	var (
		method string
		body   string
	)

	cmd := &cobra.Command{
		Use:   "api <endpoint>",
		Short: "Make an authenticated GitLab API request",
		Long: `Make a raw REST API request to GitLab.
The endpoint should start with '/' (e.g., /projects).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			hostname, token := auth.DefaultHostWithToken(cfg)
			if token == "" {
				return fmt.Errorf("authentication required: run 'glb auth login' or set GITLAB_TOKEN")
			}

			protocol := cfg.APIProtocol(hostname)
			endpoint := args[0]
			if !strings.HasPrefix(endpoint, "/") {
				endpoint = "/" + endpoint
			}

			url := fmt.Sprintf("%s://%s/api/v4%s", protocol, hostname, endpoint)

			var reqBody io.Reader
			if body != "" {
				reqBody = strings.NewReader(body)
			}

			req, err := http.NewRequestWithContext(cmd.Context(), strings.ToUpper(method), url, reqBody)
			if err != nil {
				return fmt.Errorf("creating request: %w", err)
			}

			req.Header.Set("PRIVATE-TOKEN", token)
			if body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			httpClient := &http.Client{Timeout: 30 * time.Second}
			resp, err := httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("making request: %w", err)
			}
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("reading response: %w", err)
			}

			fmt.Fprintln(f.IO.Out, string(respBody))
			return nil
		},
	}

	cmd.Flags().StringVarP(&method, "method", "X", "GET", "HTTP method")
	cmd.Flags().StringVar(&body, "body", "", "Request body (JSON)")

	return cmd
}
