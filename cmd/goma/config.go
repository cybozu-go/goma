package main

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/cybozu-go/goma"
	"github.com/cybozu-go/goma/monitor"
	"golang.org/x/net/context"
)

func loadTOML(f string) ([]*goma.MonitorDefinition, error) {
	s := &struct {
		Monitors []*goma.MonitorDefinition `toml:"monitor"`
	}{nil}
	md, err := toml.DecodeFile(f, s)
	if err != nil {
		return nil, err
	}
	if len(md.Undecoded()) > 0 {
		return nil, fmt.Errorf("undecoded keys: %v", md.Undecoded())
	}
	return s.Monitors, nil
}

func loadFile(ctx context.Context, f string) error {
	defs, err := loadTOML(f)
	if err != nil {
		return err
	}

	monitors := make([]*monitor.Monitor, 0, len(defs))
	for _, md := range defs {
		m, err := goma.CreateMonitor(md)
		if err != nil {
			return err
		}
		monitors = append(monitors, m)
	}

	for _, m := range monitors {
		// ignoring errors is safe at this point.
		monitor.Register(m)
		m.Start(ctx)
	}
	return nil
}

func loadConfigs(ctx context.Context, dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.toml"))
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := loadFile(ctx, f); err != nil {
			return err
		}
	}
	return nil
}
