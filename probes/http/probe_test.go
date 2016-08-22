package http

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cybozu-go/goma"
)

const (
	testAddress     = "localhost:13838"
	testUserAgent   = "testUserAgent"
	testHeaderName  = "X-Goma-Test"
	testHeaderValue = "gomagoma"
)

func serve(l net.Listener) {
	router := http.NewServeMux()
	router.HandleFunc("/echo/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		t := strings.Split(r.URL.Path, "/")
		w.Write([]byte(t[len(t)-1]))
	})
	router.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "500", http.StatusInternalServerError)
	})
	router.HandleFunc("/postonly", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Bad method", http.StatusBadRequest)
		}
	})
	router.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(testHeaderName) != testHeaderValue {
			http.Error(w, "Bad header", http.StatusBadRequest)
		}
	})
	router.HandleFunc("/ua", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != testUserAgent {
			http.Error(w, "Bad User Agent", http.StatusBadRequest)
		}
	})
	router.HandleFunc("/sleep/", func(w http.ResponseWriter, r *http.Request) {
		t := strings.Split(r.URL.Path, "/")
		i, err := strconv.Atoi(t[len(t)-1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		time.Sleep(time.Duration(i) * time.Second)
	})

	s := &http.Server{
		Handler: router,
	}
	s.Serve(l)
}

func TestMain(m *testing.M) {
	flag.Parse()

	l, err := net.Listen("tcp", testAddress)
	if err != nil {
		log.Fatal(err)
	}
	go serve(l)
	os.Exit(m.Run())
}

func getURL(elem ...string) string {
	return fmt.Sprintf("http://%s/%s", testAddress, strings.Join(elem, "/"))
}

func TestConstruct(t *testing.T) {
	t.Parallel()

	if _, err := construct(nil); err != goma.ErrNoKey {
		t.Error(`err != goma.ErrNoKey`)
	}

	p, err := construct(map[string]interface{}{
		"url": getURL("200"),
	})
	if err != nil {
		t.Fatal(err)
	}
	f := p.Probe(context.Background())
	if f != 0 {
		t.Error(`f != 0`)
	}
	// repeat
	f = p.Probe(context.Background())
	if f != 0 {
		t.Error(`f != 0`)
	}

	p, err = construct(map[string]interface{}{
		"url": getURL("500"),
	})
	if err != nil {
		t.Fatal(err)
	}
	f = p.Probe(context.Background())
	if !goma.FloatEquals(f, 1.0) {
		t.Error(`!goma.FloatEquals(f, 1.0)`)
	}
}

func TestHeader(t *testing.T) {
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"url": getURL("header"),
		"header": map[string]interface{}{
			testHeaderName: testHeaderValue,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	f := p.Probe(context.Background())
	if f != 0 {
		t.Error(`f != 0`)
	}

	p, err = construct(map[string]interface{}{
		"url": getURL("header"),
	})
	if err != nil {
		t.Fatal(err)
	}

	f = p.Probe(context.Background())
	if !goma.FloatEquals(f, 1.0) {
		t.Error(`!goma.FloatEquals(f, 1.0)`)
	}
}

func TestUserAgent(t *testing.T) {
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"url":   getURL("ua"),
		"agent": testUserAgent,
		"header": map[string]interface{}{
			testHeaderName: testHeaderValue,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	f := p.Probe(context.Background())
	if f != 0 {
		t.Error(`f != 0`)
	}

	p, err = construct(map[string]interface{}{
		"url":   getURL("ua"),
		"agent": testUserAgent,
	})
	if err != nil {
		t.Fatal(err)
	}

	f = p.Probe(context.Background())
	if f != 0 {
		t.Error(`f != 0`)
	}
}

func TestMethod(t *testing.T) {
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"url":    getURL("postonly"),
		"method": "POST",
	})
	if err != nil {
		t.Fatal(err)
	}

	f := p.Probe(context.Background())
	if f != 0 {
		t.Error(`f != 0`)
	}

	p, err = construct(map[string]interface{}{
		"url": getURL("postonly"),
	})
	if err != nil {
		t.Fatal(err)
	}

	f = p.Probe(context.Background())
	if !goma.FloatEquals(f, 1.0) {
		t.Error(`!goma.FloatEquals(f, 1.0)`)
	}
}

func TestProxy(t *testing.T) {
	t.Parallel()

	proxyURL := os.Getenv("GOMA_PROXY")
	if len(proxyURL) == 0 {
		t.Skip()
	}

	p, err := construct(map[string]interface{}{
		"url":   "http://example.org/",
		"proxy": proxyURL,
	})
	if err != nil {
		t.Fatal(err)
	}

	f := p.Probe(context.Background())
	if f != 0 {
		t.Error(`f != 0`)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"url":   getURL("echo", "123.45"),
		"parse": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	f := p.Probe(context.Background())
	if !goma.FloatEquals(f, 123.45) {
		t.Error(`!goma.FloatEquals(f, 123.45)`)
	}

	p, err = construct(map[string]interface{}{
		"url":    getURL("500"),
		"parse":  true,
		"errval": 100.0,
	})
	if err != nil {
		t.Fatal(err)
	}

	f = p.Probe(context.Background())
	if !goma.FloatEquals(f, 100.0) {
		t.Error(`!goma.FloatEquals(f, 100.0)`)
	}
}

func TestTimeout(t *testing.T) {
	t.Parallel()

	p, err := construct(map[string]interface{}{
		"url": getURL("sleep", "10"),
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
