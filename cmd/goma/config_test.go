package main

import (
	"os"
	"testing"
)

const content = `
[[monitor]]
[monitor.filter]
window = 2
`

func writeToTempFile(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	_, err = tmpFile.WriteString(content)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func TestLoadTOML(t *testing.T) {
	tmpFileName, err := writeToTempFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFileName)

	monitor, err := loadTOML(tmpFileName)
	if err != nil {
		t.Fatal(err)
	}

	v := monitor[0].Filter["window"]
	_, ok := v.(int64)
	if !ok {
		t.Fatalf("monitor[0].filter.window is not int64: %T", v)
	}
}
