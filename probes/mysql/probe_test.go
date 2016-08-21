package mysql

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cybozu-go/goma"
)

var (
	dsn = os.Getenv("MYSQL_DSN")
)

func TestConstruct(t *testing.T) {
	if len(dsn) == 0 {
		t.Skip("No MYSQL_DSN env")
	}
	t.Parallel()

	_, err := construct(nil)
	if err != goma.ErrNoKey {
		t.Error(`err != goma.ErrNoKey`)
	}

	_, err = construct(map[string]interface{}{
		"dsn": dsn,
	})
	if err != goma.ErrNoKey {
		t.Error(`err != goma.ErrNoKey`)
	}

	p, err := construct(map[string]interface{}{
		"dsn":   dsn,
		"query": "SELECT 1",
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	v := p.Probe(ctx)
	if !goma.FloatEquals(v, 1.0) {
		t.Error(!goma.FloatEquals(v, 1.0))
	}
}

func TestError(t *testing.T) {
	if len(dsn) == 0 {
		t.Skip("No MYSQL_DSN env")
	}
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"dsn":   dsn,
		"query": "SELECT hogenotfound()",
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	v := p.Probe(ctx)
	if v != 0 {
		t.Error(`v != 0`)
	}
}

func TestErrval(t *testing.T) {
	if len(dsn) == 0 {
		t.Skip("No MYSQL_DSN env")
	}
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"dsn":    dsn,
		"query":  "SELECT hogenotfound()",
		"errval": 123,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	v := p.Probe(ctx)
	if !goma.FloatEquals(v, 123) {
		t.Error(`!goma.FloatEquals(v, 123)`)
	}
}

func TestFloat(t *testing.T) {
	if len(dsn) == 0 {
		t.Skip("No MYSQL_DSN env")
	}
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"dsn":   dsn,
		"query": "SELECT 123.45",
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	v := p.Probe(ctx)
	if !goma.FloatEquals(v, 123.45) {
		t.Error(!goma.FloatEquals(v, 123.45))
	}
}

func TestTimeout(t *testing.T) {
	if len(dsn) == 0 {
		t.Skip("No MYSQL_DSN env")
	}
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"dsn":    dsn,
		"query":  "SELECT 100 FROM (SELECT SLEEP(10)) AS sub",
		"errval": 123.45,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	v := p.Probe(ctx)
	if !goma.FloatEquals(v, 123.45) {
		t.Error(!goma.FloatEquals(v, 123.45))
	}
}
