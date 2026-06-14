package cmd

import (
	"os"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

// legacyEndpointKeys are config keys that older CLI versions wrote into the
// user's config file. The base URL is environment config (binary default plus
// --api-url / SOVOREM_API_URL), so a value saved here would silently shadow a
// new default after an upgrade. We strip them on startup.
var legacyEndpointKeys = []string{"api_url", "frontend_url"}

// migrateConfig removes legacy endpoint keys from the active config file (if
// present) and reloads viper so its in-memory state drops them too. It is a
// best-effort, no-op-on-error cleanup: a failure here must never block the CLI.
func migrateConfig() {
	path := viper.ConfigFileUsed()
	if path == "" {
		return
	}
	if stripConfigKeys(path, legacyEndpointKeys...) {
		_ = viper.ReadInConfig()
	}
}

// stripConfigKeys removes the given top-level keys from the YAML file at path,
// preserving everything else. It returns true only if it actually rewrote the
// file.
func stripConfigKeys(path string, keys ...string) bool {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var doc map[string]any
	if err := yaml.Unmarshal(raw, &doc); err != nil || doc == nil {
		return false
	}

	removed := false
	for _, k := range keys {
		if _, ok := doc[k]; ok {
			delete(doc, k)
			removed = true
		}
	}
	if !removed {
		return false
	}

	out, err := yaml.Marshal(doc)
	if err != nil {
		return false
	}
	return os.WriteFile(path, out, 0o600) == nil
}
