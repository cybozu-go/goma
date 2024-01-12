package http

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cybozu-go/goma"
	"github.com/cybozu-go/goma/probes"
	"github.com/cybozu-go/log"
)

type probe struct {
	client *http.Client
	url    *url.URL
	method string
	header map[string]string
	parse  bool
	errval float64
}

func (p *probe) Probe(ctx context.Context) float64 {
	header := make(http.Header)
	for k, v := range p.header {
		header.Set(k, v)
	}

	req := &http.Request{
		Method:     p.method,
		URL:        p.url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     header,
		Host:       p.url.Host,
	}

	resp, err := p.client.Do(req.WithContext(ctx))
	if err != nil {
		log.Error("probe:http error", map[string]interface{}{
			"url":   p.url.String(),
			"error": err.Error(),
		})
		if p.parse {
			return p.errval
		}
		return 1.0
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		if p.parse {
			return p.errval
		}
		return 1.0
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if p.parse {
			return p.errval
		}
		return 1.0
	}

	if p.parse {
		f, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err != nil {
			log.Error("probe:http parsing failure", map[string]interface{}{
				"url":   p.url.String(),
				"error": err.Error(),
			})
			return p.errval
		}
		return f
	}
	return 0
}

func (p *probe) String() string {
	return "probe:http:" + p.url.String()
}

func construct(params map[string]interface{}) (probes.Prober, error) {
	urlStr, err := goma.GetString("url", params)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	method := "GET"
	m, err := goma.GetString("method", params)
	if err == nil {
		method = m
	} else if err != goma.ErrNoKey {
		return nil, err
	}
	header, err := goma.GetStringMap("header", params)
	switch err {
	case nil:
	case goma.ErrNoKey:
		header = make(map[string]string)
	default:
		return nil, err
	}

	switch agent, err := goma.GetString("agent", params); err {
	case nil:
		header["User-Agent"] = agent
	case goma.ErrNoKey:
		header["User-Agent"] = "goma/" + goma.Version
	default:
		return nil, err
	}

	proxy := http.ProxyFromEnvironment
	if proxyURL, err := goma.GetString("proxy", params); err == nil {
		u2, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		proxy = http.ProxyURL(u2)
	} else if err != goma.ErrNoKey {
		return nil, err
	}

	parse, err := goma.GetBool("parse", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}
	errval, err := goma.GetFloat("errval", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: proxy,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		MaxIdleConnsPerHost:   1,
		ExpectContinueTimeout: 500 * time.Millisecond,
	}
	client := &http.Client{
		Transport: transport,
	}

	return &probe{
		client: client,
		url:    u,
		method: method,
		header: header,
		parse:  parse,
		errval: errval,
	}, nil
}

func init() {
	probes.Register("http", construct)
}
