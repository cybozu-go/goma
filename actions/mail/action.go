package mail

import (
	"time"

	"github.com/cybozu-go/goma/actions"
)

type action struct {
}

func (a *action) Init(name string) error {
	return nil
}

func (a *action) Fail(name string, v float64) error {
	return nil
}

func (a *action) Recover(name string, d time.Duration) error {
	return nil
}

func (a *action) String() string {
	return "action:mail"
}

func construct(params map[string]interface{}) (actions.Actor, error) {
	return &action{}, nil
}

func init() {
	actions.Register("mail", construct)
}
