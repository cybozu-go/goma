//go:build !nacl && !plan9 && !windows
// +build !nacl,!plan9,!windows

package exec

import (
	"context"
	"testing"
	"time"

	"github.com/cybozu-go/goma"
)

func TestConstructBasic(t *testing.T) {
	t.Parallel()
	if _, err := construct(nil); err != goma.ErrNoKey {
		t.Error(`err != goma.ErrNoKey`)
	}

	if _, err := construct(map[string]interface{}{
		"command": true,
	}); err != goma.ErrInvalidType {
		t.Error(`err != goma.ErrInvalidType`)
	}
	if p, err := construct(map[string]interface{}{
		"command": "true",
	}); err != nil {
		t.Error(err)
	} else {
		if f := p.Probe(context.Background()); f != 0 {
			t.Error(`p.Probe(context.Background()) should return 0`)
		}
	}
}

func TestConstructArgs(t *testing.T) {
	t.Parallel()
	if _, err := construct(map[string]interface{}{
		"command": "echo",
		"args":    false,
	}); err != goma.ErrInvalidType {
		t.Error(`args=false should cause error`)
	}

	if p, err := construct(map[string]interface{}{
		"command": "echo",
		"args":    []interface{}{"123.45"},
		"parse":   true,
	}); err != nil {
		t.Error(err)
	} else {
		f := p.Probe(context.Background())
		if !goma.FloatEquals(f, 123.45) {
			t.Error(`!goma.FloatEquals(f, 123.45)`)
		}
	}
}

func TestProbeFalse(t *testing.T) {
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"command": "false",
	})
	if err != nil {
		t.Fatal(err)
	}

	f := p.Probe(context.Background())
	if !goma.FloatEquals(f, 1.0) {
		t.Error(`!goma.FloatEquals(f, 1.0)`)
	}
}

func TestProbeParse(t *testing.T) {
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"command": "false",
		"parse":   true,
		"errval":  3.0,
	})
	if err != nil {
		t.Fatal(err)
	}

	f := p.Probe(context.Background())
	if !goma.FloatEquals(f, 3.0) {
		t.Error(`!goma.FloatEquals(f, 3.0)`)
	}
}

func TestProbeEnv(t *testing.T) {
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"command": "sh",
		"args":    []interface{}{"-c", `echo "$GOMA_VALUE"`},
		"parse":   true,
		"env":     []interface{}{"GOMA_VALUE=123.45"},
	})
	if err != nil {
		t.Fatal(err)
	}

	f := p.Probe(context.Background())
	if !goma.FloatEquals(f, 123.45) {
		t.Error(`!goma.FloatEquals(f, 123.45)`)
	}
}

func TestProbeTimeout(t *testing.T) {
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"command": "sleep",
		"args":    []interface{}{"10"},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	f := p.Probe(ctx)
	if !goma.FloatEquals(f, 1.0) {
		t.Error(`!goma.FloatEquals(f, 1.0)`)
	}
}
