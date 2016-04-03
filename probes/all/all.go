// Package all import all probes to be compiled-in.
package all

import (
	// import all probes
	_ "github.com/cybozu-go/goma/probes/exec"
	_ "github.com/cybozu-go/goma/probes/http"
	_ "github.com/cybozu-go/goma/probes/mysql"
)
