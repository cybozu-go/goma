// Package all import all actions to be compiled-in.
package all

import (
	// import all actions
	_ "github.com/cybozu-go/goma/actions/exec"
	_ "github.com/cybozu-go/goma/actions/http"
	_ "github.com/cybozu-go/goma/actions/mail"
)
