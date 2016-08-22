package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cybozu-go/goma"
	"github.com/cybozu-go/goma/probes"
	"github.com/cybozu-go/log"
	// for driver
	_ "github.com/go-sql-driver/mysql"
)

// Obtain connection ID by "SELECT connection_id()",
// kill it by "KILL ID" to interrupt the query execution.

type probe struct {
	dsn                  string
	db                   *sql.DB
	query                string
	errval               float64
	haveMaxExecutionTime bool
}

func (p *probe) Probe(ctx context.Context) float64 {
	var err error
	var connID int64
	if p.haveMaxExecutionTime {
		var d time.Duration
		deadline, ok := ctx.Deadline()
		if ok {
			d = deadline.Sub(time.Now())
		}
		_, err = p.db.Exec("SET max_execution_time = ?", d.Nanoseconds()/1000000)
		if err != nil {
			log.Error("probe:mysql SET max_execution_time", map[string]interface{}{
				"dsn":   p.dsn,
				"error": err.Error(),
			})
			return p.errval
		}
		goto QUERY
	}
	err = p.db.QueryRow("SELECT connection_id()").Scan(&connID)
	if err != nil {
		log.Error("probe:mysql SELECT connection_id()", map[string]interface{}{
			"dsn":   p.dsn,
			"error": err.Error(),
		})
		return p.errval
	}

QUERY:
	done := make(chan float64, 1)
	go func() {
		var v float64
		err = p.db.QueryRow(p.query).Scan(&v)
		if err != nil {
			done <- p.errval
			log.Error("probe:mysql db.QueryRow", map[string]interface{}{
				"dsn":   p.dsn,
				"error": err.Error(),
			})
			return
		}
		done <- v
	}()

	select {
	case <-ctx.Done():
		if !p.haveMaxExecutionTime {
			// kill thread
			p.db.Exec("KILL ?", connID)
		}
		return p.errval
	case v := <-done:
		return v
	}
}

func (p *probe) String() string {
	return fmt.Sprintf("mysql:%s:%s", p.dsn, p.query)
}

func construct(params map[string]interface{}) (probes.Prober, error) {
	dsn, err := goma.GetString("dsn", params)
	if err != nil {
		return nil, err
	}
	query, err := goma.GetString("query", params)
	if err != nil {
		return nil, err
	}
	errval, err := goma.GetFloat("errval", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// max_execution_time is available for MySQL 5.7.8+
	_, err = db.Exec("SET max_execution_time = 10000000")
	haveMaxExecutionTime := err == nil
	if haveMaxExecutionTime {
		db.Exec("SET max_execution_time = 0")
	}

	return &probe{
		dsn:                  dsn,
		db:                   db,
		query:                query,
		errval:               errval,
		haveMaxExecutionTime: haveMaxExecutionTime,
	}, nil
}

func init() {
	probes.Register("mysql", construct)
}
