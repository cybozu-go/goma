package http

import (
	"context"
	"errors"
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
	testAddress     = "localhost:13839"
	testUserAgent   = "testUserAgent"
	testHeaderName  = "X-Goma-Test"
	testHeaderValue = "gomagoma"
	testParamName   = "test1"
	testParamValue  = "parapara"
)

func checkRequest(r *http.Request, method, event string) error {
	if r.Method != method {
		return fmt.Errorf("bad method: %s", r.Method)
	}
	if r.FormValue("monitor") != "monitor1" {
		return fmt.Errorf("bad monitor: %s", r.FormValue("monitor"))
	}
	hname, _ := os.Hostname()
	if r.FormValue("host") != hname {
		return fmt.Errorf("bad host: %s", r.FormValue("host"))
	}
	if r.FormValue("event") != event {
		return fmt.Errorf("bad event: %s", r.FormValue("event"))
	}
	if r.FormValue("version") != goma.Version {
		return fmt.Errorf("bad version: %s", r.FormValue("version"))
	}
	return nil
}

func serve(l net.Listener) {
	router := http.NewServeMux()
	router.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		if err := checkRequest(r, http.MethodGet, "init"); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	})
	router.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		if err := checkRequest(r, http.MethodGet, "fail"); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.FormValue("value") != "0.2" {
			http.Error(w, `r.FormValue("value") != "0.2"`,
				http.StatusBadRequest)
			return
		}
	})
	router.HandleFunc("/recover", func(w http.ResponseWriter, r *http.Request) {
		if err := checkRequest(r, http.MethodGet, "recover"); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.FormValue("duration") != "39" {
			http.Error(w, `r.FormValue("duration") != "39"`,
				http.StatusBadRequest)
			return
		}
	})
	router.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if err := checkRequest(r, http.MethodPost, "init"); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	})
	router.HandleFunc("/params", func(w http.ResponseWriter, r *http.Request) {
		if err := checkRequest(r, http.MethodGet, "init"); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.FormValue(testParamName) != testParamValue {
			http.Error(w, `r.FormValue(testParamName) != testParamValue`,
				http.StatusBadRequest)
			return
		}
	})
	router.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "500", http.StatusInternalServerError)
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

func makeURL(path ...string) string {
	return fmt.Sprintf("http://%s/%s", testAddress, strings.Join(path, "/"))
}

func TestConstruct(t *testing.T) {
	t.Parallel()

	_, err := construct(map[string]interface{}{
		"url_init": true,
	})
	if err != goma.ErrInvalidType {
		t.Error(`err != goma.ErrInvalidType`)
	}

	a, err := construct(nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Init("hoge"); err != nil {
		t.Error(err)
	}
}

func TestBasic(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"url_init":    makeURL("init"),
		"url_fail":    makeURL("fail"),
		"url_recover": makeURL("recover"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Init("monitor1"); err != nil {
		t.Error(err)
	}
	if err := a.Fail("monitor1", 0.2); err != nil {
		t.Error(err)
	}
	if err := a.Recover("monitor1", 39120*time.Millisecond); err != nil {
		t.Error(err)
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"url_init": makeURL("500"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Init("monitor1"); err == nil {
		t.Error("500 error is expected")
	}
}

func TestPost(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"url_init": makeURL("post"),
		"method":   http.MethodPost,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Init("monitor1"); err != nil {
		t.Error(err)
	}
}

func TestAgent(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"url_init": makeURL("ua"),
		"agent":    testUserAgent,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Init("monitor1"); err != nil {
		t.Error(err)
	}
}

func TestHeader(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"url_init": makeURL("header"),
		"header": map[string]interface{}{
			testHeaderName: testHeaderValue,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Init("monitor1"); err != nil {
		t.Error(err)
	}
}

func TestParams(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"url_init": makeURL("params"),
		"params": map[string]interface{}{
			testParamName: testParamValue,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Init("monitor1"); err != nil {
		t.Error(err)
	}
}

func TestTimeout(t *testing.T) {
	t.Parallel()

	a, err := construct(map[string]interface{}{
		"url_init": makeURL("sleep", "10"),
		"timeout":  1,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Init("monitor1"); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected %v, got %v", context.DeadlineExceeded, err)
	}
}
