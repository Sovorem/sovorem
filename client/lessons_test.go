package api

import (
	"strings"
	"testing"

	"github.com/goccy/go-json"
)

// These lock the CLI's hand-written lesson types (lessons.go) to the server's
// PascalCase /v1/lessons wire format. Unlike the auth/user types — generated
// into v1gen.go from the shared OpenAPI spec — the lesson types can't be
// generated (they double as the CLI's YAML execution model), so this is their
// drift guard: it fails if a field name/casing diverges from what the server's
// src/app/v1/lessons/[uuid]/route.ts sends. Keep in step with the lesson schemas
// in web/contracts/v1.openapi.yaml.

func TestLessonDecodesServerGetPayload(t *testing.T) {
	const payload = `{"Lesson":{"Type":"type_cli","LessonDataCLI":{"CLIData":{"BaseURLDefault":"http://localhost:8080","Steps":[{"CLICommand":{"Command":"go run main.go","Tests":[{"ExitCode":0,"StdoutContainsAll":["hello"]}]},"NoPenaltyOnFail":false}],"AllowedOperatingSystems":["linux","darwin"]}}}}`

	var l Lesson
	if err := json.Unmarshal([]byte(payload), &l); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if l.Lesson.Type != "type_cli" {
		t.Errorf("Type = %q, want type_cli", l.Lesson.Type)
	}
	if l.Lesson.LessonDataCLI == nil {
		t.Fatal("LessonDataCLI = nil, want populated")
	}
	steps := l.Lesson.LessonDataCLI.CLIData.Steps
	if len(steps) != 1 || steps[0].CLICommand == nil {
		t.Fatalf("Steps = %#v, want one CLICommand step", steps)
	}
	if steps[0].CLICommand.Command != "go run main.go" {
		t.Errorf("Command = %q, want %q", steps[0].CLICommand.Command, "go run main.go")
	}
	tests := steps[0].CLICommand.Tests
	if len(tests) != 1 || tests[0].ExitCode == nil || *tests[0].ExitCode != 0 {
		t.Errorf("Tests[0].ExitCode mismatch: %#v", tests)
	}
}

func TestLessonSubmissionEventDecodesServerSuccess(t *testing.T) {
	const payload = `{"ResultSlug":"success","StructuredErrCLI":null,"XPReward":50,"XPBreakdown":[{"Name":"Դասի ավարտ","Percent":100,"XP":50}]}`

	var ev LessonSubmissionEvent
	if err := json.Unmarshal([]byte(payload), &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ev.ResultSlug != VerificationResultSlugSuccess {
		t.Errorf("ResultSlug = %q, want success", ev.ResultSlug)
	}
	if ev.XPReward != 50 {
		t.Errorf("XPReward = %d, want 50", ev.XPReward)
	}
	if ev.StructuredErrCLI != nil {
		t.Errorf("StructuredErrCLI = %#v, want nil on success", ev.StructuredErrCLI)
	}
	if len(ev.XPBreakdown) != 1 || ev.XPBreakdown[0].XP != 50 {
		t.Errorf("XPBreakdown = %#v, want one item with XP=50", ev.XPBreakdown)
	}
}

func TestLessonSubmissionEventDecodesServerFailure(t *testing.T) {
	const payload = `{"ResultSlug":"failure","StructuredErrCLI":{"Error":"expected exit code 0, got 1","FailedStepIndex":2,"FailedTestIndex":1},"XPReward":0,"XPBreakdown":[]}`

	var ev LessonSubmissionEvent
	if err := json.Unmarshal([]byte(payload), &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ev.ResultSlug != VerificationResultSlugFailure {
		t.Errorf("ResultSlug = %q, want failure", ev.ResultSlug)
	}
	if ev.StructuredErrCLI == nil {
		t.Fatal("StructuredErrCLI = nil, want populated on failure")
	}
	// The wire key is "Error" (StructuredErrCLI.ErrorMessage has json:"Error").
	if ev.StructuredErrCLI.ErrorMessage != "expected exit code 0, got 1" {
		t.Errorf("ErrorMessage = %q", ev.StructuredErrCLI.ErrorMessage)
	}
	if ev.StructuredErrCLI.FailedStepIndex != 2 || ev.StructuredErrCLI.FailedTestIndex != 1 {
		t.Errorf("Failed indices = %d/%d, want 2/1", ev.StructuredErrCLI.FailedStepIndex, ev.StructuredErrCLI.FailedTestIndex)
	}
}

func TestCLISubmissionMarshalsPascalCase(t *testing.T) {
	body := lessonSubmissionCLI{CLIResults: []CLIStepResult{
		{CLICommandResult: &CLICommandResult{ExitCode: 0, Stdout: "hello"}},
	}}

	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	for _, want := range []string{`"CLIResults"`, `"CLICommandResult"`, `"ExitCode":0`, `"Stdout":"hello"`} {
		if !strings.Contains(got, want) {
			t.Errorf("submission JSON %s missing %s", got, want)
		}
	}
}
