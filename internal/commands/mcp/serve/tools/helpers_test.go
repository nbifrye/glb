package tools

import (
	"testing"
)

func TestSplitLabels(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"bug,feature", []string{"bug", "feature"}},
		{"bug, feature, docs", []string{"bug", "feature", "docs"}},
		{" bug , ", []string{"bug"}},
		{"", []string{}},
		{"single", []string{"single"}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := splitLabels(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("splitLabels(%q) = %v (len=%d), want %v (len=%d)", tt.input, got, len(got), tt.want, len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("splitLabels(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseIntList(t *testing.T) {
	tests := []struct {
		input string
		want  []int64
	}{
		{"1,2,3", []int64{1, 2, 3}},
		{"1, 2, 3", []int64{1, 2, 3}},
		{"42", []int64{42}},
		{"", []int64{}},
		{"1,abc,3", []int64{1, 3}},
		{" , , ", []int64{}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseIntList(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("parseIntList(%q) = %v (len=%d), want %v (len=%d)", tt.input, got, len(got), tt.want, len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseIntList(%q)[%d] = %d, want %d", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestTextResult(t *testing.T) {
	r := textResult("hello")
	if r.IsError {
		t.Error("textResult should not be an error")
	}
	if len(r.Content) != 1 {
		t.Fatalf("textResult content length = %d, want 1", len(r.Content))
	}
}

func TestErrorResult(t *testing.T) {
	r := errorResult("something went wrong")
	if !r.IsError {
		t.Error("errorResult should be an error")
	}
	if len(r.Content) != 1 {
		t.Fatalf("errorResult content length = %d, want 1", len(r.Content))
	}
}

func TestBoolPtr(t *testing.T) {
	p := boolPtr(true)
	if p == nil || !*p {
		t.Error("boolPtr(true) should return pointer to true")
	}
	p = boolPtr(false)
	if p == nil || *p {
		t.Error("boolPtr(false) should return pointer to false")
	}
}
