package average

import (
	"testing"

	"github.com/cybozu-go/goma"
)

func TestDefault(t *testing.T) {
	f, err := construct(nil)
	if err != nil {
		t.Fatal(err)
	}
	if f.Put(0) != 0 {
		t.Error("non-zero average")
	}

	f.Put(1)
	f.Put(1)
	v := f.Put(1)
	if !goma.FloatEquals(v, 0.3) {
		t.Error(`!goma.FloatEquals(v, 0.3)`)
	}
}

func TestWindow(t *testing.T) {
	_, err := construct(map[string]interface{}{
		"window": false,
	})
	if err == nil {
		t.Error(`window must be int`)
	}

	f, err := construct(map[string]interface{}{
		"window": 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	f.Put(1)
	f.Put(1)
	v := f.Put(1)
	if !goma.FloatEquals(v, 0.15) {
		t.Error(`!goma.FloatEquals(v, 0.15)`)
	}
}

func TestInit(t *testing.T) {
	_, err := construct(map[string]interface{}{
		"init": 100,
	})
	if err == nil {
		t.Error(`init must be float64`)
	}

	f, err := construct(map[string]interface{}{
		"init": 1.0,
	})
	if err != nil {
		t.Fatal(err)
	}
	f.Put(0)
	f.Put(0)
	v := f.Put(0)
	if !goma.FloatEquals(v, 0.7) {
		t.Error(`!goma.FloatEquals(v, 0.7)`)
	}
}
