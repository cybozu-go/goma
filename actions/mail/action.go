package mail

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/mail"
	"os"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/cybozu-go/goma"
	"github.com/cybozu-go/goma/actions"
	gomail "gopkg.in/gomail.v2"
)

const (
	// DefaultSubject is a text/template for Subject header.
	DefaultSubject = `alert from {{ .Monitor }} on {{ .Host }}`

	// DefaultBody is a text/template for mail message body.
	DefaultBody = `Monitor: {{ .Monitor }}
Host: {{ .Host }}
Date: {{ .Date }}
Event: {{ .Event }}
Value: {{printf "%g" .Value}}
Duration: {{ .Duration }}
Version: {{ .Version }}
`
	defaultServer = "localhost:25"
)

type tplParams struct {
	Monitor  string
	Host     string
	Date     time.Time
	Event    string
	Value    float64
	Duration int
	Version  string
}

var (
	tplSubject = template.Must(template.New("subject").Parse(DefaultSubject))
	tplBody    = template.Must(template.New("body").Parse(DefaultBody))

	headerPattern = regexp.MustCompile(`(?i)^X-[a-z0-9-]+$`)
)

type action struct {
	from      *mail.Address
	to        []*mail.Address
	initTo    []*mail.Address
	failTo    []*mail.Address
	recoverTo []*mail.Address
	subject   *template.Template
	body      *template.Template
	server    string
	user      string
	password  string
	header    map[string]string
	bcc       bool
}

func (a *action) send(params *tplParams, altTo []*mail.Address) error {
	l := len(a.to) + len(altTo)
	if l == 0 {
		return nil
	}
	to := make([]*mail.Address, 0, l)
	to = append(to, a.to...)
	to = append(to, altTo...)

	hname, err := os.Hostname()
	if err != nil {
		return err
	}
	params.Host = hname
	params.Date = time.Now()
	params.Version = goma.Version

	msg := gomail.NewMessage(gomail.SetCharset("utf-8"))
	msg.SetAddressHeader("From", a.from.Address, a.from.Name)
	sto := make([]string, 0, len(to))
	for _, t := range to {
		sto = append(sto, msg.FormatAddress(t.Address, t.Name))
	}
	rcptHeader := "To"
	if a.bcc {
		rcptHeader = "Bcc"
	}
	msg.SetHeader(rcptHeader, sto...)
	sbj := new(bytes.Buffer)
	if err := a.subject.Execute(sbj, params); err != nil {
		return err
	}
	msg.SetHeader("Subject", sbj.String())
	msg.SetDateHeader("Date", params.Date)
	for k, v := range a.header {
		msg.SetHeader(k, v)
	}

	body := new(bytes.Buffer)
	if err := a.body.Execute(body, params); err != nil {
		return err
	}
	msg.SetBody("text/plain", body.String())

	host, port, _ := net.SplitHostPort(a.server)
	nport, _ := strconv.Atoi(port)
	switch port {
	case "smtp", "mail":
		nport = 25
	case "submission":
		nport = 587
	case "urd", "ssmtp", "smtps":
		nport = 465
	}
	d := gomail.NewDialer(host, nport, a.user, a.password)
	return d.DialAndSend(msg)
}

func (a *action) Init(name string) error {
	params := &tplParams{
		Monitor: name,
		Event:   "init",
	}
	return a.send(params, a.initTo)
}

func (a *action) Fail(name string, v float64) error {
	params := &tplParams{
		Monitor: name,
		Event:   "fail",
		Value:   v,
	}
	return a.send(params, a.failTo)
}

func (a *action) Recover(name string, d time.Duration) error {
	params := &tplParams{
		Monitor:  name,
		Event:    "recover",
		Duration: int(d.Seconds()),
	}
	return a.send(params, a.recoverTo)
}

func (a *action) String() string {
	return "action:mail"
}

func getAddressList(name string, params map[string]interface{}) ([]*mail.Address, error) {
	l, err := goma.GetStringList(name, params)
	switch err {
	case nil:
		la := make([]*mail.Address, 0, len(l))
		for _, t := range l {
			a, err := mail.ParseAddress(t)
			if err != nil {
				return nil, err
			}
			la = append(la, a)
		}
		return la, nil
	case goma.ErrNoKey:
		return nil, nil
	default:
		return nil, err
	}
}

func construct(params map[string]interface{}) (actions.Actor, error) {
	fromString, err := goma.GetString("from", params)
	if err != nil {
		return nil, err
	}
	from, err := mail.ParseAddress(fromString)
	if err != nil {
		return nil, err
	}

	to, err := getAddressList("to", params)
	if err != nil {
		return nil, err
	}
	initTo, err := getAddressList("init_to", params)
	if err != nil {
		return nil, err
	}
	failTo, err := getAddressList("fail_to", params)
	if err != nil {
		return nil, err
	}
	recoverTo, err := getAddressList("recover_to", params)
	if err != nil {
		return nil, err
	}

	subject := tplSubject
	subjectString, err := goma.GetString("subject", params)
	switch err {
	case nil:
		tpl, err := template.New("subject").Parse(subjectString)
		if err != nil {
			return nil, err
		}
		if err := tpl.Execute(io.Discard, &tplParams{}); err != nil {
			return nil, err
		}
		subject = tpl
	case goma.ErrNoKey:
	default:
		return nil, err
	}

	body := tplBody
	bodyString, err := goma.GetString("body", params)
	switch err {
	case nil:
		tpl, err := template.New("body").Parse(bodyString)
		if err != nil {
			return nil, err
		}
		if err := tpl.Execute(io.Discard, &tplParams{}); err != nil {
			return nil, err
		}
		body = tpl
	case goma.ErrNoKey:
	default:
		return nil, err
	}

	server, err := goma.GetString("server", params)
	switch err {
	case nil:
		if _, _, err := net.SplitHostPort(server); err != nil {
			return nil, err
		}
	case goma.ErrNoKey:
		server = defaultServer
	default:
		return nil, err
	}

	user, err := goma.GetString("user", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}
	password, err := goma.GetString("password", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}

	header, err := goma.GetStringMap("header", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}
	for k := range header {
		if !headerPattern.MatchString(k) {
			return nil, fmt.Errorf("invalid header: %s", k)
		}
	}

	bcc, err := goma.GetBool("bcc", params)
	if err != nil && err != goma.ErrNoKey {
		return nil, err
	}

	return &action{
		from:      from,
		to:        to,
		initTo:    initTo,
		failTo:    failTo,
		recoverTo: recoverTo,
		subject:   subject,
		body:      body,
		server:    server,
		user:      user,
		password:  password,
		header:    header,
		bcc:       bcc,
	}, nil
}

func init() {
	actions.Register("mail", construct)
}
