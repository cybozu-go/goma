package exec

import (
	"context"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/cybozu-go/goma"
	"github.com/cybozu-go/goma/probes"
	"github.com/cybozu-go/log"
)

type probe struct {
	command string
	args    []string
	parse   bool
	errval  float64
	env     []string
}

func (p *probe) Probe(ctx context.Context) float64 {
	cmd := exec.CommandContext(ctx, p.command, p.args...)
	if p.env != nil {
		cmd.Env = p.env
	}

	data, err := cmd.Output()
	if err != nil {
		log.Error("probe:exec error", map[string]interface{}{
			"command": p.command,
			"args":    p.args,
			"error":   err.Error(),
		})
		if p.parse {
			return p.errval
		}
		return 1.0
	}

	if p.parse {
		f, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err != nil {
			return p.errval
		}
		return f
	}

	return 0
}

func (p *probe) String() string {
	return "probe:exec:" + p.command
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

func construct(params map[string]interface{}) (probes.Prober, error) {
	command, err := goma.GetString("command", params)
	if err != nil {
		return nil, err
	}
	args, err := goma.GetStringList("args", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}
	parse, err := goma.GetBool("parse", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}
	errval, err := goma.GetFloat("errval", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}
	env, err := goma.GetStringList("env", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}
	if env != nil {
		env = mergeEnv(env, os.Environ())
	}

	return &probe{
		command: command,
		args:    args,
		parse:   parse,
		errval:  errval,
		env:     env,
	}, nil
}

func init() {
	probes.Register("exec", construct)
}
