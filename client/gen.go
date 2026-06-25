package api

// The CLI's pure /v1 wire types — auth + user (OtpLoginRequest, CliTokens,
// CurrentUser, LogoutResponse, Error) — are GENERATED from an OpenAPI contract
// into v1gen.go. Do not hand-edit them.
//
// v1.openapi.yaml here is a vendored copy of the GENERATED portion of
// web/contracts/v1.openapi.yaml (the source of truth in the sovorem-am repo).
// That spec ALSO documents the lesson delivery + grading wire (LessonResponse,
// CLISubmission, LessonSubmissionEvent, …), but those schemas are deliberately
// NOT vendored or generated here: the CLI's lesson types live by hand in
// lessons.go because they double as its local YAML execution model (see
// cmd/localtest.go) — carrying yaml tags, methods and constants a generated type
// can't. lessons_test.go locks those hand-written types to the same wire shapes.
// So: when the spec's auth/user schemas change, sync them here and regenerate;
// when its lesson schemas change, update lessons.go + lessons_test.go to match.
//
// Regenerate v1gen.go after editing the spec:
//
//	go generate ./...

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.0 -generate types -package api -o v1gen.go v1.openapi.yaml
