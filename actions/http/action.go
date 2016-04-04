package http

import (
	"time"

	"github.com/cybozu-go/goma/actions"
)

type action struct {
}

func (a *action) Init(name string) {
}

func (a *action) Fail(name string, v float64) {
}

func (a *action) Recover(name string, d time.Duration) {
}

func (a *action) String() string {
	return "action:http"
}

func construct(params map[string]interface{}) (actions.Actor, error) {
	return &action{}, nil
}

func init() {
	actions.Register("http", construct)
}
