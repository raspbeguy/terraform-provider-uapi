//go:build tools

// Package tools pins the codegen tools used by `make docs` so they are tracked
// in go.mod and reproducible. It is never compiled into the provider binary.
package tools

import (
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
