package exec

import (
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/cybozu-go/goma"
	"github.com/cybozu-go/goma/probes"
	"github.com/cybozu-go/log"
	"golang.org/x/net/context"
)

type probe struct {
	command string
	args    []string
	parse   bool
	errval  float64
	env     []string
}

func (p *probe) createCmd() *exec.Cmd {
	cmd := exec.Command(p.command, p.args...)
	cmd.Dir = "/"
	if p.env != nil {
		cmd.Env = p.env
	}
	return cmd
}

type retType struct {
	f   float64
	err error
}

func (p *probe) Probe(ctx context.Context) (float64, error) {
	cmd := p.createCmd()
	ch := make(chan retType, 1)

	go func() {
		data, err := cmd.Output()
		if p.parse {
			f, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
			if err != nil {
				ch <- retType{p.errval, err}
			} else {
				ch <- retType{f, nil}
			}
			return
		}
		if err != nil {
			ch <- retType{1.0, err}
			return
		}
		ch <- retType{0, nil}
	}()

	select {
	case ret := <-ch:
		return ret.f, ret.err
	case <-ctx.Done():
		cmd.Process.Kill()
		if log.Enabled(log.LvDebug) {
			log.Debug("probe:exec killed", map[string]interface{}{
				"_command": p.command,
			})
		}
		if p.parse {
			return p.errval, ctx.Err()
		}
		return 1.0, ctx.Err()
	}
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
