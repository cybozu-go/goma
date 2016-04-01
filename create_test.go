package goma

import (
	"encoding/json"
	"testing"

	"github.com/BurntSushi/toml"
)

const (
	testFile = "sample.toml"
)

func testMonitor1(t *testing.T, m *MonitorDefinition) {
	if m.Name != "monitor1" {
		t.Error(`m.Name != "monitor1"`)
	}
	if m.Interval != 10 {
		t.Error(`m.Interval != 10`)
	}
	if m.Timeout != 1 {
		t.Error(`m.Timeout != 1`)
	}
	if !FloatEquals(m.Min, 0) {
		t.Error(`!FloatEquals(m.Min, 0)`)
	}
	if !FloatEquals(m.Max, 0.3) {
		t.Error(`!FloatEquals(m.Max, 0.3)`)
	}
	if pt, err := getType(m.Probe); err != nil {
		t.Error(err)
	} else if pt != "exec" {
		t.Error(`pt != "exec"`)
	}
	if pcmd, err := GetStringList("command", m.Probe); err != nil {
		t.Error(err)
	} else if len(pcmd) != 1 {
		t.Error(`len(pcmd) != 1`)
	} else if pcmd[0] != "/some/probe/cmd" {
		t.Error(`pcmd[0] != "/some/probe/cmd"`)
	}
}

func TestSample(t *testing.T) {
	t.Parallel()

	d := &struct {
		Monitors []*MonitorDefinition `toml:"monitor"`
	}{}
	_, err := toml.DecodeFile(testFile, d)
	if err != nil {
		t.Fatal(err)
	}

	if len(d.Monitors) != 2 {
		t.Fatal("len(d.Monitors) != 2, ", len(d.Monitors))
	}

	m1 := d.Monitors[0]
	testMonitor1(t, m1)
	data, err := json.Marshal(m1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	var jm1 MonitorDefinition
	err = json.Unmarshal(data, &jm1)
	if err != nil {
		t.Fatal(err)
	}
	testMonitor1(t, &jm1)
}
