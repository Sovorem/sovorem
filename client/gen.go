package api

// The CLI's /v1 wire types are generated from an OpenAPI contract that is the
// single source of truth shared with the web backend. v1.openapi.yaml here is a
// vendored copy of web/contracts/v1.openapi.yaml in the sovorem-am repo — keep
// the two in sync when the contract changes.
//
// Regenerate v1gen.go after editing the spec:
//
//	go generate ./...

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.0 -generate types -package api -o v1gen.go v1.openapi.yaml
