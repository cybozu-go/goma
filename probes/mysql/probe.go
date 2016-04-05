package mysql

import (
	"github.com/cybozu-go/goma/probes"
	"golang.org/x/net/context"
)

// Obtain connection ID by "SELECT connection_id()",
// kill it by "KILL ID" to interrupt the query execution.

type probe struct {
}

func (p *probe) Probe(ctx context.Context) float64 {
	return 0
}

func (p *probe) String() string {
	return "mysql"
}

func construct(params map[string]interface{}) (probes.Prober, error) {
	return &probe{}, nil
}

func init() {
	probes.Register("mysql", construct)
}
