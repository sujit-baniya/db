//go:build oracle
// +build oracle

package migrations

import (
	_ "github.com/mattn/go-oci8"
)

func init() {
	dialects["oci8"] = OracleDialect{}
}
