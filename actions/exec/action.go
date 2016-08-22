package exec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/cybozu-go/goma"
	"github.com/cybozu-go/goma/actions"
	"github.com/cybozu-go/log"
)

const (
	eventInit    = "init"
	eventFail    = "fail"
	eventRecover = "recover"

	envMonitor  = "GOMA_MONITOR"
	envEvent    = "GOMA_EVENT"
	envValue    = "GOMA_VALUE"
	envDuration = "GOMA_DURATION"
	envVersion  = "GOMA_VERSION"
)

type action struct {
	command string
	args    []string
	env     []string
	timeout time.Duration
	debug   bool
}

func mergeEnv(env, bgenv []string) (merged []string) {
	m := make(map[string]string)
	for _, e := range bgenv {
		m[strings.SplitN(e, "=", 2)[0]] = e
	}
	for _, e := range env {
		m[strings.SplitN(e, "=", 2)[0]] = e
	}
	for _, v := range m {
		merged = append(merged, v)
	}
	sort.Strings(merged)
	return
}

func (a *action) run(env []string) error {
	var cmd *exec.Cmd
	if a.timeout == 0 {
		cmd = exec.Command(a.command, a.args...)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, a.command, a.args...)
	}
	cmd.Dir = "/"
	cmd.Env = mergeEnv(env, a.env)

	if a.debug {
		out, err := cmd.CombinedOutput()
		log.Error("action:exec debug", map[string]interface{}{
			"output": out,
			"error":  err.Error(),
		})
		return err
	}
	return cmd.Run()
}

func (a *action) Init(name string) error {
	env := []string{
		fmt.Sprintf("%s=%s", envMonitor, name),
		fmt.Sprintf("%s=%s", envVersion, goma.Version),
		fmt.Sprintf("%s=%s", envEvent, eventInit),
	}
	return a.run(env)
}

func (a *action) Fail(name string, v float64) error {
	env := []string{
		fmt.Sprintf("%s=%s", envMonitor, name),
		fmt.Sprintf("%s=%s", envVersion, goma.Version),
		fmt.Sprintf("%s=%s", envEvent, eventFail),
		fmt.Sprintf("%s=%g", envValue, v), // %g suppresses trailing zeroes.
	}
	return a.run(env)
}

func (a *action) Recover(name string, d time.Duration) error {
	env := []string{
		fmt.Sprintf("%s=%s", envMonitor, name),
		fmt.Sprintf("%s=%s", envVersion, goma.Version),
		fmt.Sprintf("%s=%s", envEvent, eventRecover),
		fmt.Sprintf("%s=%d", envDuration, int(d.Seconds())),
	}
	return a.run(env)
}

func (a *action) String() string {
	return "action:exec:" + a.command
}

func construct(params map[string]interface{}) (actions.Actor, error) {
	command, err := goma.GetString("command", params)
	if err != nil {
		return nil, err
	}
	args, err := goma.GetStringList("args", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}
	env, err := goma.GetStringList("env", params)
	switch err {
	case nil:
		env = mergeEnv(env, os.Environ())
	case goma.ErrNoKey:
		env = os.Environ()
	default:
		return nil, err
	}
	timeout, err := goma.GetInt("timeout", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}
	debug, err := goma.GetBool("debug", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}

	return &action{
		command: command,
		args:    args,
		env:     env,
		timeout: time.Duration(timeout) * time.Second,
		debug:   debug,
	}, nil
}

func init() {
	actions.Register("exec", construct)
}
