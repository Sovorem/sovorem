package api

import (
	"testing"

	"github.com/goccy/go-json"
)

// These lock the CLI's wire format to the server's PascalCase /v1 contract.
// The types under test are generated from v1.openapi.yaml (the single source of
// truth), so a contract change that broke this would show up at generate time.

func TestOtpLoginRequestMarshalsPascalCase(t *testing.T) {
	b, err := json.Marshal(OtpLoginRequest{Otp: "abc123"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(b) != `{"Otp":"abc123"}` {
		t.Errorf("OtpLoginRequest marshaled to %s, want {\"Otp\":\"abc123\"}", b)
	}
}

func TestCliTokensDecodesServerPayload(t *testing.T) {
	var r CliTokens
	if err := json.Unmarshal([]byte(`{"AccessToken":"sov_at_x","RefreshToken":"sov_rt_y"}`), &r); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if r.AccessToken != "sov_at_x" || r.RefreshToken != "sov_rt_y" {
		t.Errorf("decoded %+v, want AccessToken/RefreshToken populated", r)
	}
}

func TestCurrentUserDecodesServerPayload(t *testing.T) {
	var u CurrentUser
	if err := json.Unmarshal([]byte(`{"Handle":"nazani"}`), &u); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if u.Handle != "nazani" {
		t.Errorf("decoded Handle=%q, want nazani", u.Handle)
	}
}
