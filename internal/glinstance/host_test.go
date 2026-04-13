package glinstance

import "testing"

func TestDefault(t *testing.T) {
	if got := Default(); got != "gitlab.com" {
		t.Errorf("Default() = %q, want %q", got, "gitlab.com")
	}
}

func TestNormalizeHostname(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"gitlab.com", "gitlab.com"},
		{"GITLAB.COM", "gitlab.com"},
		{"https://gitlab.com", "gitlab.com"},
		{"http://gitlab.com", "gitlab.com"},
		{"https://gitlab.com/", "gitlab.com"},
		{"http://GITLAB.COM/", "gitlab.com"},
		{"my-gitlab.example.com", "my-gitlab.example.com"},
		{"https://my-gitlab.example.com/", "my-gitlab.example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := NormalizeHostname(tt.input); got != tt.want {
				t.Errorf("NormalizeHostname(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAPIEndpoint(t *testing.T) {
	tests := []struct {
		hostname string
		protocol string
		want     string
	}{
		{"gitlab.com", "https", "https://gitlab.com/api/v4"},
		{"gitlab.com", "", "https://gitlab.com/api/v4"},
		{"my.host.com", "http", "http://my.host.com/api/v4"},
		{"https://gitlab.com/", "https", "https://gitlab.com/api/v4"},
	}
	for _, tt := range tests {
		t.Run(tt.hostname+"_"+tt.protocol, func(t *testing.T) {
			if got := APIEndpoint(tt.hostname, tt.protocol); got != tt.want {
				t.Errorf("APIEndpoint(%q, %q) = %q, want %q", tt.hostname, tt.protocol, got, tt.want)
			}
		})
	}
}
