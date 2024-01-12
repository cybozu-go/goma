//go:build !nacl && !plan9 && !windows
// +build !nacl,!plan9,!windows

package exec

import (
	"testing"
	"time"

	"github.com/cybozu-go/goma"
)

func TestConstruct(t *testing.T) {
	t.Parallel()

	_, err := construct(nil)
	if err != goma.ErrNoKey {
		t.Error(`err != goma.ErrNoKey`)
	}

	a, err := construct(map[string]interface{}{
		"command": "sh",
		"args": []interface{}{"-u", "-c", `
echo GOMA_MONITOR=$GOMA_MONITOR
echo GOMA_VERSION=$GOMA_VERSION
if [ "$GOMA_EVENT" != "init" ]; then exit 1; fi
`},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Init("monitor1"); err != nil {
		t.Error(err)
	}
}

func TestFail(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"command": "sh",
		"args": []interface{}{"-u", "-c", `
echo GOMA_VALUE=$GOMA_VALUE
if [ "$GOMA_EVENT" != "fail" ]; then exit 1; fi
if [ "$GOMA_VALUE" != "0.1" ]; then exit 1; fi
`},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Fail("monitor1", 0.1); err != nil {
		t.Error(err)
	}
}

func TestRecover(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"command": "sh",
		"args": []interface{}{"-u", "-c", `
echo GOMA_DURATION=$GOMA_DURATION
if [ "$GOMA_EVENT" != "recover" ]; then exit 1; fi
if [ "$GOMA_DURATION" != "39" ]; then exit 1; fi
`},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Recover("monitor1", 39280*time.Millisecond); err != nil {
		t.Error(err)
	}
}

func TestEnv(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"command": "sh",
		"args": []interface{}{"-u", "-c", `
echo TEST_ENV1=$TEST_ENV1
if [ "$TEST_ENV1" != "test1" ]; then exit 1; fi
`},
		"env": []interface{}{"TEST_ENV1=test1"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Fail("monitor1", 0); err != nil {
		t.Error(err)
	}
}

func TestTimeout(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"command": "sleep",
		"args":    []interface{}{"10"},
		"timeout": 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Init("monitor1"); err == nil {
		t.Error("err must not be nil")
	}
}
