package spec

import _ "embed"

// OpenAPI 스펙 원본을 바이너리에 포함
//
//go:embed api/openapi.yaml
var OpenAPI []byte
