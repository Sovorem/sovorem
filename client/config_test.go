package api

import "testing"

func TestNormalizeBaseURL(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"https://sovorem.am", "https://sovorem.am", false},
		{"https://sovorem.am/", "https://sovorem.am", false},
		{"http://localhost:3000", "http://localhost:3000", false},
		{"http://localhost:3000/", "http://localhost:3000", false},
		{"  https://stage.sovorem.am  ", "https://stage.sovorem.am", false},
		{"localhost:3000", "", true},
		{"ftp://example.com", "", true},
		{"", "", true},
	}
	for _, c := range cases {
		got, err := NormalizeBaseURL(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("NormalizeBaseURL(%q): expected error, got %q", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("NormalizeBaseURL(%q): unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("NormalizeBaseURL(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestBaseURLPrecedence(t *testing.T) {
	t.Cleanup(func() { baseOverride = "" })

	// Default: no flag, no env -> the compiled-in production origin.
	baseOverride = ""
	if got := APIBaseURL(); got != DefaultBaseURL {
		t.Errorf("APIBaseURL default = %q, want %q", got, DefaultBaseURL)
	}
	if got := FrontendBaseURL(); got != DefaultBaseURL {
		t.Errorf("FrontendBaseURL default = %q, want %q", got, DefaultBaseURL)
	}

	// Env beats the default and is trimmed; the two env vars are independent.
	t.Setenv("SOVOREM_API_URL", "https://api.example.com/")
	t.Setenv("SOVOREM_FRONTEND_URL", "https://example.com/")
	if got := APIBaseURL(); got != "https://api.example.com" {
		t.Errorf("APIBaseURL with env = %q", got)
	}
	if got := FrontendBaseURL(); got != "https://example.com" {
		t.Errorf("FrontendBaseURL with env = %q", got)
	}

	// The --api-url override beats env for BOTH resolvers and is normalized.
	if err := SetBaseOverride("http://localhost:3000/"); err != nil {
		t.Fatalf("SetBaseOverride: %v", err)
	}
	if got := APIBaseURL(); got != "http://localhost:3000" {
		t.Errorf("APIBaseURL with override = %q, want http://localhost:3000", got)
	}
	if got := FrontendBaseURL(); got != "http://localhost:3000" {
		t.Errorf("FrontendBaseURL with override = %q, want http://localhost:3000", got)
	}

	// An empty value is a no-op (leaves the existing override in place).
	if err := SetBaseOverride(""); err != nil {
		t.Fatalf("SetBaseOverride(\"\"): %v", err)
	}
	if got := APIBaseURL(); got != "http://localhost:3000" {
		t.Errorf("APIBaseURL after empty SetBaseOverride = %q", got)
	}

	// An invalid value errors and does not change the override.
	if err := SetBaseOverride("localhost:3000"); err == nil {
		t.Error("SetBaseOverride(\"localhost:3000\"): expected error")
	}
	if got := APIBaseURL(); got != "http://localhost:3000" {
		t.Errorf("APIBaseURL after invalid SetBaseOverride = %q (override should be unchanged)", got)
	}
}
