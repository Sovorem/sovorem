package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestStripConfigKeysRemovesLegacyEndpointsAndKeepsTheRest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	original := []byte(`api_url: https://api.sovorem.am
frontend_url: https://sovorem.am
access_token: secret-access
refresh_token: secret-refresh
last_refresh: 1700000000
color:
  red: "1"
  green: "2"
`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if !stripConfigKeys(path, legacyEndpointKeys...) {
		t.Fatal("stripConfigKeys returned false, expected it to rewrite the file")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	if _, ok := doc["api_url"]; ok {
		t.Error("api_url was not removed")
	}
	if _, ok := doc["frontend_url"]; ok {
		t.Error("frontend_url was not removed")
	}
	// User data must be preserved untouched.
	if doc["access_token"] != "secret-access" {
		t.Errorf("access_token = %v, want secret-access", doc["access_token"])
	}
	if doc["refresh_token"] != "secret-refresh" {
		t.Errorf("refresh_token = %v, want secret-refresh", doc["refresh_token"])
	}
	if _, ok := doc["color"]; !ok {
		t.Error("color map was dropped")
	}
}

func TestStripConfigKeysNoOpWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	original := []byte("access_token: secret\nlast_refresh: 42\n")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if stripConfigKeys(path, legacyEndpointKeys...) {
		t.Error("stripConfigKeys returned true, expected no-op when keys absent")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(raw) != string(original) {
		t.Errorf("file was rewritten on a no-op:\n got: %q\nwant: %q", raw, original)
	}
}

func TestStripConfigKeysMissingFileIsSafe(t *testing.T) {
	if stripConfigKeys(filepath.Join(t.TempDir(), "does-not-exist.yaml"), legacyEndpointKeys...) {
		t.Error("expected false for a missing file")
	}
}
