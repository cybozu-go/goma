package http

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cybozu-go/goma"
	"github.com/cybozu-go/goma/actions"
)

const (
	defaultTimeout = 30
)

var (
	client = &http.Client{}
)

type action struct {
	urlInit    *url.URL
	urlFail    *url.URL
	urlRecover *url.URL
	method     string
	header     map[string]string
	params     map[string]string
	timeout    time.Duration
}

func processResponse(u *url.URL, resp *http.Response) error {
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if 200 <= resp.StatusCode && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("action:http:%s %s", u.String(), resp.Status)
}

func (a *action) request(u *url.URL, params map[string]string) error {
	tu := *u
	values := tu.Query()
	for k, v := range a.params {
		values.Set(k, v)
	}
	for k, v := range params {
		values.Set(k, v)
	}
	hname, err := os.Hostname()
	if err != nil {
		return err
	}
	values.Set("host", hname)
	values.Set("version", goma.Version)
	data := values.Encode()

	header := make(http.Header)
	for k, v := range a.header {
		header.Set(k, v)
	}

	var body io.ReadCloser
	var length int64
	if a.method == http.MethodGet {
		tu.RawQuery = data
	} else {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		length = int64(len(data))
		body = ioutil.NopCloser(strings.NewReader(data))
	}
	req := &http.Request{
		Method:        a.method,
		URL:           &tu,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		Body:          body,
		ContentLength: length,
		Host:          u.Host,
	}

	if a.timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	return processResponse(u, resp)
}

func (a *action) Init(name string) error {
	if a.urlInit == nil {
		return nil
	}
	params := make(map[string]string)
	for k, v := range a.params {
		params[k] = v
	}
	params["monitor"] = name
	params["event"] = "init"
	return a.request(a.urlInit, params)
}

func (a *action) Fail(name string, v float64) error {
	if a.urlFail == nil {
		return nil
	}
	params := make(map[string]string)
	for k, v := range a.params {
		params[k] = v
	}
	params["monitor"] = name
	params["event"] = "fail"
	params["value"] = fmt.Sprintf("%g", v) // %g suppresses trailing zeroes.
	return a.request(a.urlFail, params)
}

func (a *action) Recover(name string, d time.Duration) error {
	if a.urlRecover == nil {
		return nil
	}
	params := make(map[string]string)
	for k, v := range a.params {
		params[k] = v
	}
	params["monitor"] = name
	params["event"] = "recover"
	params["duration"] = strconv.Itoa(int(d.Seconds()))
	return a.request(a.urlRecover, params)
}

func (a *action) String() string {
	return fmt.Sprintf("action:http:%s:%s:%s",
		a.urlInit, a.urlFail, a.urlRecover)
}

func construct(params map[string]interface{}) (actions.Actor, error) {
	var uI, uF, uR *url.URL
	urlInit, err := goma.GetString("url_init", params)
	switch err {
	case nil:
		uI, err = url.Parse(urlInit)
		if err != nil {
			return nil, err
		}
	case goma.ErrNoKey:
	default:
		return nil, err
	}
	urlFail, err := goma.GetString("url_fail", params)
	switch err {
	case nil:
		uF, err = url.Parse(urlFail)
		if err != nil {
			return nil, err
		}
	case goma.ErrNoKey:
	default:
		return nil, err
	}
	urlRecover, err := goma.GetString("url_recover", params)
	switch err {
	case nil:
		uR, err = url.Parse(urlRecover)
		if err != nil {
			return nil, err
		}
	case goma.ErrNoKey:
	default:
		return nil, err
	}

	method, err := goma.GetString("method", params)
	switch err {
	case nil:
	case goma.ErrNoKey:
		method = http.MethodGet
	default:
		return nil, err
	}
	agent, err := goma.GetString("agent", params)
	switch err {
	case nil:
	case goma.ErrNoKey:
		agent = "goma/" + goma.Version
	default:
		return nil, err
	}
	header, err := goma.GetStringMap("header", params)
	switch err {
	case nil:
	case goma.ErrNoKey:
		header = map[string]string{"User-Agent": agent}
	default:
		return nil, err
	}
	formParams, err := goma.GetStringMap("params", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}
	timeout, err := goma.GetInt("timeout", params)
	switch err {
	case nil:
	case goma.ErrNoKey:
		timeout = defaultTimeout
	default:
		return nil, err
	}

	return &action{
		urlInit:    uI,
		urlFail:    uF,
		urlRecover: uR,
		method:     method,
		header:     header,
		params:     formParams,
		timeout:    time.Duration(timeout) * time.Second,
	}, nil
}

func init() {
	actions.Register("http", construct)
}
