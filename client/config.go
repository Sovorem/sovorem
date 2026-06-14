package api

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// DefaultBaseURL is the production origin the CLI talks to. It is compiled into
// the binary on purpose: the endpoint is environment configuration, not user
// data, so it lives in code (shipped via releases) and is never persisted to a
// user's config file where a stale value could shadow it after an upgrade.
const DefaultBaseURL = "https://sovorem.am"

// baseOverride, when non-empty, takes precedence over everything for the
// duration of a single CLI invocation. It is populated from the --api-url flag
// (see cmd.initConfig) and is deliberately kept out of viper so it is never
// written to the user's config file.
var baseOverride string

// SetBaseOverride validates raw and, when non-empty, makes it the base origin
// for both API requests and the browser login URL. An empty string is a no-op,
// so callers can pass an unset flag value safely.
func SetBaseOverride(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	normalized, err := NormalizeBaseURL(raw)
	if err != nil {
		return err
	}
	baseOverride = normalized
	return nil
}

// APIBaseURL returns the origin used for /v1/* API requests, without a trailing
// slash. Precedence: --api-url flag > SOVOREM_API_URL env > built-in default.
func APIBaseURL() string {
	return resolveBaseURL("SOVOREM_API_URL")
}

// FrontendBaseURL returns the origin used for the browser /cli/login page,
// without a trailing slash. Precedence: --api-url flag > SOVOREM_FRONTEND_URL
// env > built-in default.
func FrontendBaseURL() string {
	return resolveBaseURL("SOVOREM_FRONTEND_URL")
}

func resolveBaseURL(envVar string) string {
	if baseOverride != "" {
		return baseOverride
	}
	if v := strings.TrimRight(strings.TrimSpace(os.Getenv(envVar)), "/"); v != "" {
		return v
	}
	return DefaultBaseURL
}

// NormalizeBaseURL validates that raw is an absolute http(s) origin and returns
// it without a trailing slash. http is intentionally allowed so local dev
// servers (e.g. http://localhost:3000) work.
func NormalizeBaseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("անվավեր URL %q: %w", raw, err)
	}
	// url.Parse reads "localhost:3000" as scheme "localhost" with an empty host,
	// so require both a real scheme and a host.
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("անվավեր base URL %q. ներառիր և՛ scheme-ը (http/https), և՛ host-ը, օր․ http://localhost:3000", raw)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("base URL-ի scheme-ը պետք ա լինի http կամ https, ոչ թե %q", u.Scheme)
	}
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}
