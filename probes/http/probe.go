package http

import (
	"github.com/cybozu-go/goma/probes"
	"golang.org/x/net/context"
)

type probe struct {
}

func (p *probe) Probe(ctx context.Context) (float64, error) {
	return 0, nil
}

func (p *probe) String() string {
	return "http"
}

func construct(params map[string]interface{}) (probes.Prober, error) {
	return &probe{}, nil
}

func init() {
	probes.Register("http", construct)
}
