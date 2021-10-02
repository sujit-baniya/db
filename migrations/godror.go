//go:build godror
// +build godror

// godror is another oracle driver
// repo: https://github.com/godror/godror
//
// godror package don't cofigure pkg config on your machine,
// it mean that we don't need to specify oracle office client
// at compile process and just config oracle client at runtime.
package migrations

import (
	_ "github.com/godror/godror"
)

func init() {
	dialects["godror"] = OracleDialect{}
}
